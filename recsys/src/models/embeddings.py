from pydantic import BaseModel


class SparseVector(BaseModel):
    indices: list[int]
    values: list[float]


class VectorRepresentation(BaseModel):
    dense: list[float] | None = None
    sparse: SparseVector | None = None
