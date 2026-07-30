import logging
from datetime import UTC, datetime

from aioredis.client import Pipeline, Redis
from aioredis.exceptions import ConnectionError, RedisError
from pydantic import ValidationError

from src.core.exc import InternalError
from src.domain.schemas import ApplicationMessage

ZSET_NAME = "messages"


class MessageRedisSaver:
    def __init__(self, client: Redis) -> None:
        self.client: Redis = client

        self.logger = logging.getLogger("MessageRedisSaver")

    async def save_messages(self, messages: list[ApplicationMessage]) -> None:
        async with self.client.pipeline() as pipeline:
            try:
                names = await self.save_or_update_messages(pipeline, messages)
                await self.save_messages_in_zset(pipeline, names)
                await pipeline.execute()
            except (RedisError, ConnectionError) as e:
                self.logger.critical(f"save_messages: can't save messages, details={e}")
                raise InternalError("can't save messages in redis saver")

    async def save_or_update_messages(self, pipeline: Pipeline, messages: list[ApplicationMessage]) -> list[str]:
        names = [""] * len(messages)

        for i, message in enumerate(messages):
            names[i] = generate_name(message)

            if await self.client.exists(names[i]):
                pipeline.hincrby(name=names[i], key="amount", amount=message.amount)
            else:
                pipeline.hset(
                    name=names[i],
                    mapping={
                        "vacancy_id": message.vacancy_id,
                        "amount": message.amount,
                    },
                )

        return names

    async def save_messages_in_zset(self, pipeline: Pipeline, names: list[str]) -> None:
        save_time = datetime.now(UTC).timestamp()

        for name in names:
            if not await self.client.zscore(ZSET_NAME, name):
                pipeline.zadd(ZSET_NAME, mapping={name: save_time})

    async def get_last_messages(self, size: int) -> list[ApplicationMessage] | None:
        current_time = datetime.now(UTC).timestamp()

        names = await self.client.zrangebyscore(name=ZSET_NAME, min=0, max=current_time, start=0, num=size)
        if not names:
            return None

        async with self.client.pipeline() as pipeline:
            [pipeline.hgetall(name) for name in names]

            try:
                return [
                    ApplicationMessage(vacancy_id=int(data["vacancy_id"]), amount=int(data["amount"]))
                    for data in await pipeline.execute()
                ]
            except (ValidationError, KeyError, ValueError) as e:
                self.logger.critical(f"get_last_message: invalid keys in redis, details={e}")
                return None

            except (RedisError, ConnectionError) as e:
                self.logger.critical(f"save_messages: can't get last messages, details={e}")
                raise InternalError("can't get last messages from redis saver")

    async def delete_last_messages(self, size: int) -> None:
        current_time = datetime.now(UTC).timestamp()

        names = await self.client.zrangebyscore(ZSET_NAME, min=0, max=current_time, start=0, num=size)
        async with self.client.pipeline() as pipeline:
            try:
                pipeline.zrem(ZSET_NAME, *names)
                await pipeline.execute()
            except (RedisError, ConnectionError) as e:
                self.logger.critical(f"save_messages: can't save messages, details={e}")
                raise InternalError("can't save messages in redis saver")


def generate_name(message: ApplicationMessage) -> str:
    return f"message:vacancy:{message.vacancy_id}"
