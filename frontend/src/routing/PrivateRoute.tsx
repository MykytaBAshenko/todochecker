// components/routes/PrivateRoute.tsx
import { useSelector } from 'react-redux';
import { Navigate, Outlet } from 'react-router-dom';
import { RootState } from '../contexts/store';
import Layout from './layout/Layout';
// import Layout2 from '../Layout2'; // Your authenticated layout
// import Loader from '../Loader'; // Optional: you can make this a spinner or full screen

const PrivateRoute = () => {
  const { user, isLoading, isTokenExist } = useSelector((state: RootState) => state.auth);

  if (isLoading || (isTokenExist && !user)) {
    return <div>Loading</div>;
  }

  if (!user) {
    return <Navigate to="/sign-in" />;
  }

  return (
    <Layout>
      <Outlet />
    </Layout>
  );
};

export default PrivateRoute;
