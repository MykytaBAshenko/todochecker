import { createSlice, PayloadAction } from '@reduxjs/toolkit';
import { loginUser, registerUser, validateUserToken, logoutUser } from '../actions/authActions';

type AuthState = {
  user: any | null;
  isLoading: boolean;
  error: string | null;
  isTokenExist: boolean;
};

const initialState: AuthState = {
  user: null,
  isLoading: false,
  error: null,
  isTokenExist: !!localStorage.getItem("token")
};

const authSlice = createSlice({
  name: 'auth',
  initialState,
  reducers: {
    setUser(state, action: PayloadAction<any>) {
      state.user = action.payload;
      state.isTokenExist = !!localStorage.getItem("token")
    },
    clearUser(state) {
      state.user = null;
      state.isTokenExist = !!localStorage.getItem("token")
    },
    clearError(state) {
      state.error = null;
      state.isTokenExist = !!localStorage.getItem("token")
    }
  },
  extraReducers: (builder) => {
    // Login
    builder.addCase(loginUser.pending, (state) => {
      state.isLoading = true;
      state.error = null;
    });
    builder.addCase(loginUser.fulfilled, (state, action) => {
      state.isLoading = false;
      state.user = action.payload;
      state.error = null;
    });
    builder.addCase(loginUser.rejected, (state, action) => {
      state.isLoading = false;
      state.error = action.payload as string;
    });

    // Register
    builder.addCase(registerUser.pending, (state) => {
      state.isLoading = true;
      state.error = null;
    });
    builder.addCase(registerUser.fulfilled, (state, action) => {
      state.isLoading = false;
      state.user = action.payload;
      state.error = null;
    });
    builder.addCase(registerUser.rejected, (state, action) => {
      state.isLoading = false;
      state.error = action.payload as string;
    });

    // Validate Token
    builder.addCase(validateUserToken.pending, (state) => {
      state.isLoading = true;
      state.isTokenExist = !!localStorage.getItem("token");
    });
    builder.addCase(validateUserToken.fulfilled, (state, action) => {
      state.isLoading = false;
      state.user = action.payload;
      state.isTokenExist = !!localStorage.getItem("token");
    });
    builder.addCase(validateUserToken.rejected, (state) => {
      state.isLoading = false;
      state.user = null;
      state.isTokenExist = !!localStorage.getItem("token");
    });

    // Logout
    builder.addCase(logoutUser.fulfilled, (state) => {
      state.user = null;
      state.isTokenExist = !!localStorage.getItem("token");
    });
  },
});

export const { setUser, clearUser, clearError } = authSlice.actions;
export default authSlice.reducer;