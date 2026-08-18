from hypothesis import given
from hypothesis import strategies as st

from src.domain.schemas import ApplicationMessage
from src.infrastructure.service.application.application_service import create_update_applications_fields


@given(
    message_ids=st.lists(st.uuids(version=4), min_size=10, max_size=10),
    vacancy_ids=st.lists(st.integers(min_value=0), min_size=10, max_size=10),
)
def test_create_update_applications_fields(message_ids, vacancy_ids) -> None:
    assert len(message_ids) == len(vacancy_ids)

    apps = {
        ApplicationMessage(message_id=message_id, vacancy_id=vacancy_id)
        for message_id, vacancy_id in zip(message_ids, vacancy_ids)
    }

    vacancy_ids = {app.vacancy_id for app in apps}
    fields = create_update_applications_fields(apps)

    for field in fields:
        assert field["vacancy_id"] in vacancy_ids


@given(message_id=st.uuids(version=4), vacancy_ids=st.lists(st.integers(min_value=0), min_size=5, max_size=20))
def test_duplicates_messages(message_id, vacancy_ids) -> None:
    apps = {ApplicationMessage(message_id=message_id, vacancy_id=vacancy_id) for vacancy_id in vacancy_ids}

    fields = create_update_applications_fields(apps)
    assert len(fields) == 1
