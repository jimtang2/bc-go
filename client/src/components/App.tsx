import React from 'react';
import { RouterProvider, createBrowserRouter, Navigate } from 'react-router-dom';
import Layout from "./Layout"
import MatchesPage from './matches';
import AlphaPage from './alpha';
// import SettingsPage from './settings';

const router = createBrowserRouter([{
  path: '/',
  element: <Layout />,
  children: [
    { index: true, element: <MatchesPage /> },
    { path: 'alpha', element: <AlphaPage />} ,
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