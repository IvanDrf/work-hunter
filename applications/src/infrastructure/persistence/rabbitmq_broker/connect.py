from aio_pika import ExchangeType, connect_robust
from aio_pika.abc import AbstractChannel, AbstractExchange, AbstractRobustConnection

from src.core.config import RabbitMQConfig


async def connect_to_rabbitmq(config: RabbitMQConfig) -> AbstractRobustConnection:
    return await connect_robust(
        host=config.rabbitmq_host,
        port=config.rabbitmq_port,
        login=config.rabbitmq_user,
        password=config.rabbitmq_password,
    )


async def declare_channel(conn: AbstractRobustConnection) -> AbstractChannel:
    return await conn.channel()


async def declare_exchange(config: RabbitMQConfig, chan: AbstractChannel) -> AbstractExchange:
    return await chan.declare_exchange(name=config.rabbitmq_exchange, durable=True, type=ExchangeType.DIRECT)
