import { useSelector } from "react-redux";
import { NavLink } from "react-router-dom";
import { RootState } from "../contexts/store";
import {UserControll} from "./UserControll";

const PublicHeader = () => {
  const user = useSelector((state: RootState) => state.auth.user);

  return (
    <header className="header">
      <div className="header-container">
        <span className="header-title">
          <NavLink to="/">ToDoChecker</NavLink>
        </span>
     
          {user ? (
              <UserControll />
          ) : (
            <span className="header-buttons">
              <NavLink
                to="/register"
                className="header-button"
              >
                Register
              </NavLink>
              <NavLink
                to="/sign-in"
                className="header-button"
              >
                Sign In
              </NavLink>
            </span>

          )}
      </div>
    </header>
  );
};

export default PublicHeader;