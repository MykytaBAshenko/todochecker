import { useSelector } from "react-redux";
import { NavLink } from "react-router-dom";
import { RootState } from "../contexts/store";
import { useAppContext } from "../contexts/AppContext";
import { SignOutButton } from "./SignOutButton";
import CanvasAvatar  from "./CanvasAvatar";
import { useState } from "react";
import { MdSms } from "react-icons/md";

const UserControll = () => {
  const user = useSelector((state: RootState) => state.auth.user);
  const [openMenu, setOpenMenu] = useState(false);


  return (
    <div className="user-controll">
      <NavLink to="/messanger"><MdSms/></NavLink>
      <div
        onClick={() => setOpenMenu(!openMenu)}
        className="user-controll-display"
      >
        <div className="user-controll-nickname">
          {user.nickname}
        </div>
        <div className="user-controll-avatar">
          <CanvasAvatar avatar={user.avatar}/>
        </div>
        {openMenu ? (
          <div className="user-controll-menu">
            <NavLink to="/settings">Settings</NavLink>
            <SignOutButton />
          </div>
        ) : null}
      </div>
    </div>
  );
};

export { UserControll };
