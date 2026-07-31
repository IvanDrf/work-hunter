import logging
from asyncio import wait_for

from aio_pika.abc import AbstractChannel, AbstractIncomingMessage, AbstractQueue, AbstractRobustConnection
from pydantic import ValidationError

from src.api.rabbitmq.dependencies import IApplicationService
from src.core.exc import ArgumentError, InternalError
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

        self.logger = logging.getLogger("RabbitMQConsumer")

    async def stop(self) -> None:
        await self.chan.close()
        await self.conn.close()

        await self.application_service.stop()

    async def start_consuming(self) -> None:
        self.logger.info("Starting consumer")
        async with self.queue.iterator() as it:
            async for message in it:
                self.logger.info(f"message=got, {message.message_id=}")
                await self._process_message(message)

    async def _process_message(self, message: AbstractIncomingMessage) -> None:
        try:
            application = ApplicationMessage.model_validate_json(message.body)
            await wait_for(self.application_service.increase_vacancy_applications(application), timeout=self.service_timeout)
            await message.ack()
            self.logger.info(f"message=handled, {message.message_id=}")

        except ValidationError as e:
            self.logger.error(f"invalid incoming message, message_id={message.message_id}, details={e}")
            await message.nack(requeue=False)

        except ArgumentError as e:
            self.logger.error(f"invalid argument in message, details={e}")
            await message.nack(requeue=False)

        except InternalError as e:
            self.logger.critical(f"can't update vacancy applications, details={e}")
            await message.nack(requeue=True)

        except TimeoutError as e:
            self.logger.error(f"can't update vacancy applications, details={e}")
            await message.nack(requeue=True)

        except OSError as e:
            self.logger.critical(f"can't update vacancy applications, details={e}")
            await message.nack(requeue=True)
