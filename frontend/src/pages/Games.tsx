
import "./Games.scss";
import { useLocation, matchPath, NavLink } from "react-router-dom";

const Games = () => {
  
  return (
    <div className="games">
          <NavLink
                    to="/games/tik-tak-toe"
                  >
                    tik-tak-toe
                  </NavLink>
    </div>
  );
};

export default Games;
