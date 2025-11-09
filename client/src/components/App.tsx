import React from 'react';
import { RouterProvider, createBrowserRouter, Navigate } from 'react-router-dom';
import Layout from "./Layout"
import SpreadsTable from './spreads';
import AlphaTable from './alpha';
// import SettingsPage from './settings';

const router = createBrowserRouter([{
  path: '/',
  element: <Layout />,
  children: [
    { index: true, element: <SpreadsTable /> },
    { path: 'alpha', element: <AlphaTable />} ,
    // { path: 'settings', element: <SettingsPage /> },
    { path: '*', element: <Navigate to="/" replace /> },
  ],
}]);

export default () => {
  return (
    <React.StrictMode>
      <div className="App">
        <RouterProvider router={router} />
      </div>
    </React.StrictMode>
  );
};;