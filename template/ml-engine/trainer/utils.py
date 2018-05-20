"""Utilities."""

import os
from tensorflow.python.keras.datasets import mnist
from tensorflow.python.lib.io import file_io

MNIST_MODEL = 'mnist.hdf5'


def load_data():
  """
  Load train and test data.
  """
  x_train = None
  y_train = None
  x_test = None
  y_test = None
  (x_train, y_train), (x_test, y_test) = mnist.load_data()
  return (x_train, y_train), (x_test, y_test)


def save_model(model, job_dir):
  # Unhappy hack to work around h5py not being able to write to GCS.
  # Force snapshots and saves to local filesystem, then copy them over to GCS.
  if job_dir.startswith("gs://"):
    model.save(MNIST_MODEL)
    copy_file_to_gcs(job_dir, MNIST_MODEL)
  else:
    model.save(os.path.join(job_dir, MNIST_MODEL))

# h5py workaround: copy local models over to GCS if the job_dir is GCS.
def copy_file_to_gcs(job_dir, file_path):
  with file_io.FileIO(file_path, mode='r') as input_f:
    with file_io.FileIO(os.path.join(job_dir, file_path), mode='w+') as output_f:
        output_f.write(input_f.read())
