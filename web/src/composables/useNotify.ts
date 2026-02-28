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

  return { notImplemented };
}
