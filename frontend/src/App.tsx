import {
  BrowserRouter as Router,
  Route,
  Routes,
  Navigate,
} from "react-router-dom";
import PrivateRoute from "./routing/PrivateRoute";
import PublicRoute from "./routing/PublicRoute";
import Register from "./pages/Register";
import SignIn from "./pages/SignIn";
import Booking from "./pages/Booking";
import Messanger from "./pages/Messanger";
import Games from "./pages/Games";
import Home from "./pages/Home";
import TikTakToe from "./games/tic-tac-toe/TikTakToe";

const App = () => {
  return (
    <Router>
      <Routes>
        {/* Public Routes Group */}
        <Route element={<PublicRoute />}>
          <Route path="/" element={<Home />} />
          <Route path="/register" element={<Register />} />
          <Route path="/sign-in" element={<SignIn />} />
        </Route>

        {/* Private Routes Group */}
        <Route element={<PrivateRoute />}>
          <Route path="/dashboard" element={<Booking />} />
          <Route path="/messanger" element={<Messanger />} />
          <Route path="/group" element={<Messanger />} />
          <Route path="/games" element={<Games />} />
          <Route path="/games/tik-tak-toe" element={<TikTakToe />} />

        </Route>

        {/* Fallback */}
        <Route path="*" element={<Navigate to="/" />} />
      </Routes>
    </Router>
  );
};

export default App;
