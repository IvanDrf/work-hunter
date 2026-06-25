import logging
from asyncio import Semaphore, create_task

from aio_pika import RobustChannel, RobustConnection, RobustQueue
from aio_pika.abc import AbstractIncomingMessage
from pydantic import ValidationError

from src.api.rabbitmq.dependencies import IApplicationService
from src.core.exc import AccessError, ArgumentError, InternalError
from src.domain.schemas import ApplicationMessage


class Consumer:
    def __init__(
        self,
        conn: RobustConnection,
        chan: RobustChannel,
        queue: RobustQueue,
        application_service: IApplicationService,
        parallel_messages_amount: int,
    ) -> None:
        self.conn: RobustConnection = conn
        self.chan: RobustChannel = chan
        self.queue: RobustQueue = queue

        self.application_service: IApplicationService = application_service

        self.sem: Semaphore = Semaphore(parallel_messages_amount)

    async def start_consuming(self) -> None:
        async with self.queue.iterator() as it:
            async for message in it:
                create_task(self._process_message(message))

    async def _process_message(self, message: AbstractIncomingMessage) -> None:
        try:
            async with self.sem:
                application = ApplicationMessage.model_validate_json(message.body)
                await self.application_service.increase_vacancy_applications(application)

        except ValidationError as e:
            logging.error(f"invalid incoming message, message_id={message.message_id}, details={e}")
            await message.nack(requeue=False)

        except AccessError as e:
            logging.error(f"invalid user role for job applying, details={e}")
            await message.nack(requeue=False)

        except ArgumentError as e:
            logging.error(f"invalid argument in message, details={e}")
            await message.nack(requeue=False)

        except InternalError as e:
            logging.error(f"can't update vacancy applications, details={e}")
            await message.nack(requeue=True)
