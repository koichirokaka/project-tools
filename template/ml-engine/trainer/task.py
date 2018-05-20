"""Main program on Cloud ML"""

import argparse
import os
import model
import utils
import tensorflow.python.keras as keras
from tensorflow.python.keras import backend as K

# input image dimensions
img_rows, img_cols = 28, 28


def main(train_files,
         eval_files,
         job_dir,
         epochs,
         batch_size,
         num_classes):
  """Main program"""
  (x_train, y_train), (x_test, y_test) = utils.load_data(train_files, eval_files)
  if K.image_data_format() == 'channes_first':
    x_train = x_train.reshape(x_train.shape[0], 1, img_rows, img_cols)
    x_test = x_test.reshape(x_test.shape[0], 1, img_rows, img_cols)
    input_shape = (1, img_rows, img_cols)
  else:
    x_train = x_train.reshape(x_train.shape[0], img_rows, img_cols, 1)
    x_test = x_test.reshape(x_test.shape[0], img_rows, img_cols, 1)
    input_shape = (img_rows, img_cols, 1)

  x_train = x_train.astype('float32')
  x_test = x_test.astype('float32')
  x_train /= 255
  x_test /= 255

  y_train = keras.utils.to_categorical(y_train, num_classes)
  y_test = keras.utils.to_categorical(y_test, num_classes)
  mnist_model = model.build_model(input_shape, num_classes)
  mnist_model.fit(x_train, y_train,
                  batch_size=batch_size,
                  epochs=epochs,
                  verbose=1,
                  validation_data=(x_test, y_test))

  # Convert the Keras model to TensorFlow SavedModel
  model.to_savedmodel(mnist_model, os.path.join(job_dir, 'export'))


if __name__ == '__main__':
  parser = argparse.ArgumentParser()
  parser.add_argument('--train-files',
                      required=True,
                      type=str,
                      help='Training files')
  parser.add_argument('--eval-file',
                      required=True,
                      type=str,
                      help='Evaluation files')
  parser.add_argument('--job-dir',
                      required=True,
                      type=str,
                      help='Write checkpoints and export model')
  parser.add_argument('--epochs',
                      defalt=12,
                      required=True,
                      type=int,
                      help='Epoch')
  parser.add_argument('--batch-size',
                      default=128,
                      required=True,
                      type=int,
                      help='Batch size')
  parser.add_argument('--num-classes',
                      default=10,
                      required=True,
                      type=int,
                      help='Number of class')

  args = parser.parse_args()

  main(**args.__dict__)
