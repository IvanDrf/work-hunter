from sqlalchemy import select
from sqlalchemy.dialects.postgresql import insert
from sqlalchemy.exc import DBAPIError, SQLAlchemyError

from src.core.exc import InternalError
from src.domain.models import TagORM
from src.infrastructure.persistence.postgresql_repo.unit_of_work import UnitOfWork
from src.utils.catch_error import catch_raise_error


class TagRepo:
    @catch_raise_error((SQLAlchemyError, DBAPIError), InternalError, "critical", "can't add tags in database")
    async def add_tags(self, uof: UnitOfWork, tags: list[str]) -> list[TagORM]:
        if not tags:
            return []

        tags_values = [{"tag": tag} for tag in tags]
        await uof.session.execute(insert(TagORM).values(tags_values).on_conflict_do_nothing(index_elements=["tag"]))

        res = await uof.session.execute(select(TagORM).where(TagORM.tag.in_([t["tag"] for t in tags_values])))
        return list(res.scalars().all())
