import logging
from datetime import UTC, datetime

from redis import ConnectionError, RedisError
from redis.asyncio import Redis
from redis.asyncio.client import Pipeline

from src.core.exc import InternalError
from src.domain.schemas import ApplicationMessage


class MessageRedisSaver:
    def __init__(self, client: Redis) -> None:
        self.client: Redis = client

        self.logger = logging.getLogger("MessageRedisSaver")

    async def save_messages(self, messages: list[ApplicationMessage]) -> None:
        pipeline = self.client.pipeline(transaction=True)

        try:
            names = await save_or_update_messages(pipeline, messages)
            await save_messages_in_zset(pipeline, names)
            await pipeline.execute()
        except (RedisError, ConnectionError) as e:
            self.logger.critical(f"save_messages: can't save messages, details={e}")
            raise InternalError("can't save messages in redis saver")

        finally:
            await pipeline.aclose()


async def save_or_update_messages(pipeline: Pipeline, messages: list[ApplicationMessage]) -> list[str]:
    names = [""] * len(messages)

    for i, message in enumerate(messages):
        names[i] = generate_name(message)

        if await pipeline.exists(names[i]):
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


async def save_messages_in_zset(pipeline: Pipeline, names: list[str]) -> None:
    SET_NAME = "messages"
    save_time = datetime.now(UTC).timestamp()

    for name in names:
        if not await pipeline.zscore(SET_NAME, name):
            await pipeline.zadd(SET_NAME, mapping={name: save_time})


def generate_name(message: ApplicationMessage) -> str:
    return f"message:vacancy:{message.vacancy_id}"
