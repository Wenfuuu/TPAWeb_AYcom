import { writable } from "svelte/store";
import { nanoid } from "nanoid";

export type ToastType = "success" | "error" | "warning" | "info";
export type ToastPosition =
  | "top-right"
  | "top-left"
  | "bottom-right"
  | "bottom-left"
  | "top-center"
  | "bottom-center";

export interface Toast {
  id: string;
  type: ToastType;
  message: string;
  duration: number;
  position: ToastPosition;
}

export interface ToastOptions {
  type?: ToastType;
  duration?: number;
  position?: ToastPosition;
}

const defaultOptions: ToastOptions = {
  type: "info",
  duration: 3000,
  position: "top-right",
};

function createToastStore() {
  const { subscribe, update } = writable<Toast[]>([]);

  function addToast(message: string, options?: ToastOptions) {
    const id = nanoid();
    const toast: Toast = {
      id,
      message,
      type: options?.type || (defaultOptions.type as ToastType),
      duration: options?.duration ?? (defaultOptions.duration as number),
      position: options?.position || (defaultOptions.position as ToastPosition),
    };

    update((toasts) => [...toasts, toast]);
    return id;
  }

  function removeToast(id: string) {
    update((toasts) => toasts.filter((toast) => toast.id !== id));
  }

  function success(message: string, options?: Omit<ToastOptions, "type">) {
    return addToast(message, { ...options, type: "success" });
  }

  function error(message: string, options?: Omit<ToastOptions, "type">) {
    return addToast(message, { ...options, type: "error" });
  }

  function warning(message: string, options?: Omit<ToastOptions, "type">) {
    return addToast(message, { ...options, type: "warning" });
  }

  function info(message: string, options?: Omit<ToastOptions, "type">) {
    return addToast(message, { ...options, type: "info" });
  }

  return {
    subscribe,
    add: addToast,
    remove: removeToast,
    success,
    error,
    warning,
    info,
  };
}

export const toastStore = createToastStore();
