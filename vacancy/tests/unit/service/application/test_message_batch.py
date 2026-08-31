from hypothesis import given
from hypothesis import strategies as st
from pytest import raises

from src.domain.schemas import ApplicationMessage
from src.infrastructure.service.application.application_service import MessageBatch


@given(
    message_id=st.uuids(version=4),
    vacancy_id=st.integers(min_value=0),
    message_ids=st.lists(st.uuids(version=4), min_size=5, max_size=5),
    vacancy_ids=st.lists(st.integers(min_value=0), min_size=5, max_size=5),
)
def test_message(message_id, vacancy_id, message_ids, vacancy_ids) -> None:
    _test_one_size_batch(message_id, vacancy_id)
    _test_normal_size_batch(message_ids, vacancy_ids)


@given(valid_size=st.integers(min_value=1), invalid_size=st.integers(max_value=0))
def test_invalid_size_batch(valid_size, invalid_size) -> None:
    batch = MessageBatch(size=valid_size)
    assert batch.size == valid_size

    with raises(ValueError):
        batch = MessageBatch(size=invalid_size)


def _test_one_size_batch(message_id, vacancy_id) -> None:
    application = ApplicationMessage(message_id=message_id, vacancy_id=vacancy_id)

    batch = MessageBatch(size=1)
    batch.add_application(application)

    with raises(OverflowError):
        batch.add_application(application)

    app = batch.get_applications()
    assert not batch.get_applications()
    assert len(app) == 1

    app = app.pop()

    assert app.message_id == application.message_id
    assert app.vacancy_id == application.vacancy_id


def _test_normal_size_batch(message_ids, vacancy_ids) -> None:
    assert len(message_ids) == len(vacancy_ids)

    size = len(message_ids) // 2

    batch = MessageBatch(size)
    apps = [
        ApplicationMessage(message_id=message_id, vacancy_id=vacancy_id)
        for message_id, vacancy_id in zip(message_ids, vacancy_ids)
    ]

    assert not batch.get_applications()

    for i, app in enumerate(apps, start=1):
        if i > size:
            with raises(OverflowError):
                batch.add_application(app)

        else:
            batch.add_application(app)

    saved_apps = batch.get_applications()
    for app in apps[:size]:
        assert app in saved_apps

    for app in apps[size:]:
        assert app not in saved_apps

    assert not batch.get_applications()
