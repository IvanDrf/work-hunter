import logging

from aio_pika import Message
from aio_pika.abc import AbstractChannel, AbstractExchange, AbstractRobustConnection
from aio_pika.exceptions import ConnectionClosed, DeliveryError, PublishError

from src.core.exc import InternalError
from src.domain.schemas import ApplicationMessage, Messages


class ApplicationsProducer:
    def __init__(
        self, conn: AbstractRobustConnection, chan: AbstractChannel, exchange: AbstractExchange, routing_key: str
    ) -> None:
        self.conn = conn
        self.chan = chan
        self.exchange = exchange
        self.routing_key = routing_key

        self.logger = logging.getLogger("ApplicationsProducer")

    async def publish_application(self, messages: list[ApplicationMessage]) -> None:
        try:
            await self.exchange.publish(
                message=Message(body=Messages.dump_json(messages)),
                routing_key=self.routing_key,
            )
        except (ConnectionClosed, PublishError, DeliveryError) as e:
            self.logger.critical(f"publish_application: can't send applications, details={e}")
            raise InternalError("can't send applications to vacany service")
