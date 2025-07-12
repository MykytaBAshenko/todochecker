import { useSelector } from "react-redux";
import { RootState } from "../contexts/store";
import { useState } from "react";
import { useLocation, matchPath, NavLink } from "react-router-dom";
import "./Sidebar.scss";
import { FaBars, FaTimes, FaArrowLeft } from "react-icons/fa";

const DefaultSidebar = () => {
    const [isOpen, setIsOpen] = useState(false);
    return (
      <div className="default-sidebar">
        <div className="sidebar-controll">
            {isOpen ? <FaTimes onClick={() => setIsOpen(!isOpen)}/> : <FaBars  onClick={() => setIsOpen(!isOpen)}/>}
        </div>
        <nav>
          <NavLink
            to="/dashboard"
            className={({ isActive }) => (isActive ? "active" : "")}
          >
            <FaBars /> { isOpen ? "Dashboard" : null}
          </NavLink>
          <NavLink
            to="/messanger"
            className={({ isActive }) => (isActive ? "active" : "")}
          >
                        <FaBars /> { isOpen ? "Messanger" : null}
          </NavLink>
                    <NavLink
            to="/games"
            className={({ isActive }) => (isActive ? "active" : "")}
          >
                        <FaBars /> { isOpen ? "Games" : null}
          </NavLink>
        </nav>
      </div>
    );
  };
  


  const GroupSidebar = () => {
    const [isOpen, setIsOpen] = useState(false);
    return (
      <div className="group-sidebar">
        <div className="sidebar-controll">
            <NavLink
                to="/dashboard"
                className={({ isActive }) => (isActive ? "active" : "")}
            >
                <FaArrowLeft />
            </NavLink>
            {isOpen ? <FaTimes onClick={() => setIsOpen(!isOpen)}/> : <FaBars  onClick={() => setIsOpen(!isOpen)}/>}

        </div>
        <nav>
          <NavLink
            to="/group"
            className={({ isActive }) => (isActive ? "active" : "")}
            end
          >
                                    <FaBars /> { isOpen ? "Group" : null}
          </NavLink>
          <NavLink
            to="/group/123"
            className={({ isActive }) => (isActive ? "active" : "")}
          >
                                                <FaBars /> { isOpen ? "Group data" : null}

          </NavLink>
        </nav>
      </div>
    );
  };

const Sidebar = () => {
  const user = useSelector((state: RootState) => state.auth.user);
  const location = useLocation();
  
  const isGroupRoute = () => {
    let paths_array = ["/group/:id", "/group"]
    
    for (let i = 0; i < paths_array.length; i++) {
        if (location.pathname === paths_array[i] || matchPath(paths_array[i], location.pathname)) {
            return true;
        }
    }   
    return false;
  }
//   const isGroupRoute = matchPath("/group/:id", location.pathname) || location.pathname === "/group";

  return (
    <aside className="sidebar">
      {isGroupRoute() ? <GroupSidebar /> : <DefaultSidebar />}
    </aside>
  );
};

export default Sidebar;