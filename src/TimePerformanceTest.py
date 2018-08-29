# coding: utf-8

from time import perf_counter


def Tester(func):
    """时间测量
    """
    def get_time_spent(*args, **kwargs):
        start = perf_counter()
        func(*args, **kwargs)
        end = perf_counter()

        print(f"{func.__repr__} Spent Time: {(end-start):.6f}s")
    return get_time_spent


@Tester
def hello():
    print('hello, world!')


if __name__ == '__main__':
    hello()
