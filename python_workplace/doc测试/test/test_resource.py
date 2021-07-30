import sys
import os
import pytest

sys.path.insert(0, "../")
print(os.getcwd())



def test_add():
    print(os.getcwd())
    assert 6 == 6