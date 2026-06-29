from pydantic import BaseModel


class ValidateMetroSchema(BaseModel):
    city: str
    metro: str
