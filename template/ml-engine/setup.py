from setuptools import find_packages
from setuptools import setup

REQUIRED_PACKAGES = ['h5py==2.7.0']

setup(
    name='%s',
    version='0.1',
    install_requires=REQUIRED_PACKAGES,
    packages=find_packages(),
    include_package_data=True,
    description=''
)
