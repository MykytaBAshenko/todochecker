import Header from "../../components/Header";
import Sidebar from "../../components/Sidebar";

import { useSelector } from "react-redux";
import { RootState } from "../../contexts/store";

interface Props {
  children: React.ReactNode;
}

const PublicLayout = ({ children }: Props) => {
  const user = useSelector((state: RootState) => state.auth.user);

  return (
    <div className="main-body">
      { user &&
        <Sidebar />
      }
      <div className="main-body-content">
        <Header />
        <div className="main-body-content-container">{children}</div>
      </div>
    </div>
  );
};

export default PublicLayout;
