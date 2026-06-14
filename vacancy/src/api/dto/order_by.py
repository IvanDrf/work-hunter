from pkg.vacancy_api.vacancy_pb2 import OrderBy as PKGOrderBy

from src.core.exc import ArgumentError
from src.domain.types.enums import OrderBy


def order_by_dto(request) -> OrderBy:
    if not request.HasField("order_by"):
        return OrderBy.CREATED_AT

    order_by = {
        PKGOrderBy.DATE: OrderBy.CREATED_AT,
        PKGOrderBy.VIEWS: OrderBy.VIEWS,
        PKGOrderBy.APPLICATIONS: OrderBy.APPLICATIONS_COUNT,
    }

    if request.order_by not in order_by:
        raise ArgumentError(f"invalid order_by in request for search, {request.order_by=}")

    return order_by[request.order_by]
