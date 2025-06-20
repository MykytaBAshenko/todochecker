// useToast.ts
import { useDispatch } from "react-redux";
import { showToast, hideToast } from "../contexts/slices/toastSlice";

type ToastProps = {
    message: string;
    type: "SUCCESS" | "ERROR"
  };

export const useToast = () => {
  const dispatch = useDispatch();

  const triggerToast = ({message, type = "SUCCESS"}: ToastProps) => {
    dispatch(showToast({ message, type }));

    setTimeout(() => {
      dispatch(hideToast());
    }, 5000);
  };

  return { triggerToast };
};
