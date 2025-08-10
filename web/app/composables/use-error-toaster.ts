import { toast } from "vue-sonner";

type FetchError = {
  message: string;
  data?: { error: string };
};

export const useErrorToaster = () => {
  return (error: FetchError) => {
    const message = error.data?.error || error.message;

    const id = toast.error(message, {
      action: {
        label: "Dismiss",
        onClick: () => toast.dismiss(id),
      },
    });

    return id;
  };
};
