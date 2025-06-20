import { createSlice, PayloadAction } from '@reduxjs/toolkit';

type ToastMessage = {
  message: string;
  type: 'SUCCESS' | 'ERROR';
} | null;

const toastSlice = createSlice({
  name: 'toast',
  initialState: null as ToastMessage,
  reducers: {
    showToast(_, action: PayloadAction<ToastMessage>) {
      return action.payload;
    },
    hideToast() {
      return null;
    },
  },
});

export const { showToast, hideToast } = toastSlice.actions;
export default toastSlice.reducer;
