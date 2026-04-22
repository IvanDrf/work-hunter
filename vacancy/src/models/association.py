from sqlalchemy import Column, ForeignKey, Table

from src.models.base import Base


association_table = Table(
    'vacancies_to_tags',
    Base.metadata,
    Column('vacancy_id', ForeignKey('vacancies.vacancy_id'), primary_key=True),
    Column('tag_id', ForeignKey('tags.tag_id'), primary_key=True)
)
