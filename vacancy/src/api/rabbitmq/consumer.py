import logging
from asyncio import wait_for

from aio_pika.abc import AbstractChannel, AbstractIncomingMessage, AbstractQueue, AbstractRobustConnection
from pydantic import ValidationError

from src.api.rabbitmq.dependencies import IApplicationService
from src.core.exc import AccessError, ArgumentError, InternalError
from src.domain.schemas import ApplicationMessage


class RabbitMQConsumer:
    def __init__(
        self,
        conn: AbstractRobustConnection,
        chan: AbstractChannel,
        queue: AbstractQueue,
        application_service: IApplicationService,
        service_timeout: float,
    ) -> None:
        self.conn: AbstractRobustConnection = conn
        self.chan: AbstractChannel = chan
        self.queue: AbstractQueue = queue

        self.application_service: IApplicationService = application_service
        self.service_timeout: float = service_timeout

    async def start_consuming(self) -> None:
        logging.info("Starting consumer")
        async with self.queue.iterator() as it:
            async for message in it:
                logging.info(f"message=got, {message.message_id=}")
                await self._process_message(message)

    async def _process_message(self, message: AbstractIncomingMessage) -> None:
        try:
            application = ApplicationMessage.model_validate_json(message.body)
            await wait_for(self.application_service.increase_vacancy_applications(application), timeout=self.service_timeout)
            await message.ack()
            logging.info(f"message=handled, {message.message_id=}")

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
            logging.critical(f"can't update vacancy applications, details={e}")
            await message.nack(requeue=True)

        except TimeoutError as e:
            logging.error(f"can't update vacancy applications, details={e}")
            await message.nack(requeue=True)

        except OSError as e:
            logging.critical(f"can't update vacancy applications, details={e}")
            await message.nack(requeue=True)
