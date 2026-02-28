import { Notify } from 'quasar';

export function useNotify() {
  const notImplemented = () => {
    Notify.create({
      color: 'tertiary',
      message: 'Not yet implemented',
      actions: [
        {
          icon: 'close',
          color: 'white',
          round: true,
          handler: () => {},
        },
      ],
    });
  };

  const error = (message: string) => {
    Notify.create({
      color: 'error',
      message: message,
      actions: [
        {
          icon: 'close',
          color: 'white',
          round: true,
          handler: () => {},
        },
      ],
    });
  };

  return { error, notImplemented };
}
