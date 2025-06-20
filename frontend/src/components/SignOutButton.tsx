
import { useDispatch } from "react-redux";
import { useNavigate } from "react-router-dom";
import { clearUser } from "../contexts/slices/authSlice";
import { showToast } from "../contexts/slices/toastSlice";
import * as apiClient from "../api-client";
import { useAppContext } from "../contexts/AppContext";

const SignOutButton = () => {
  const dispatch = useDispatch();
  // Keep for backward compatibility

  const handleSignOut = async () => {
    await apiClient.signOut();
    
    dispatch(clearUser());
    dispatch(showToast({ message: "Signed Out Successfully", type: "SUCCESS" }));
  };
  return (
    <button
      onClick={handleSignOut}
      className="header-button"
    >
      Sign Out
    </button>
  );
};

export {SignOutButton}
