from pydantic import Field
from src.core.config.base import BaseConfig


class RabbitMQConfig(BaseConfig):
    rabbitmq_host: str = Field(default="localhost", validation_alias="RABBITMQ_HOST")
    rabbitmq_port: int = Field(default=5672, validation_alias="RABBITMQ_PORT")

    rabbitmq_user: str = Field(default="user", validation_alias="RABBITMQ_USER")
    rabbitmq_password: str = Field(default="password", validation_alias="RABBITMQ_PASSWORD")

    rabbitmq_exchange: str = Field(default="exchange", validation_alias="RABBITMQ_EXCHANGE")
    rabbitmq_routing_key: str = Field(default="key", validation_alias="RABBITMQ_ROUTING_KEY")
    rabbitmq_producer_queue: str = Field(default="queue", validation_alias="RABBITMQ_PRODUCER_QUEUE")
