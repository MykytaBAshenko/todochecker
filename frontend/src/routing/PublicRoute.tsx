import Layout from './layout/Layout';
import { Outlet } from 'react-router-dom';

const PublicRoute = () => (
  <Layout>
    <Outlet />
  </Layout>
);

export default PublicRoute;