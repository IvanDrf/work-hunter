from src.domain.types.types import SingleTon


def test_singleton() -> None:
    class A(metaclass=SingleTon):
        pass

    class B(metaclass=SingleTon):
        pass

    a1, a2 = A(), A()
    assert a1 is a2

    b1, b2 = B(), B()
    assert b1 is b2

    assert (a1 is b1) is False
    assert (a1 is b2) is False
    assert (a2 is b1) is False
    assert (a2 is b2) is False
