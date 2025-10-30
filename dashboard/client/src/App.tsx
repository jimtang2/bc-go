// webhook$$bc-go;dashboard/client/src/App.tsx;grok$$
import React from 'react';
import { RouterProvider, createBrowserRouter, Navigate } from 'react-router-dom';
import RootLayout from './Dashboard';
import AlphaTable from './components/table-alpha';
import HistoryTable from './components/table-history';
import Settings from './components/form-settings';

const router = createBrowserRouter([
  {
    path: '/',
    element: <RootLayout />,
    children: [
      { index: true, element: <AlphaTable /> },
      { path: 'alpha', element: <AlphaTable /> },
      { path: 'settings', element: <Settings /> },
      { path: 'history', element: <HistoryTable /> },
      { path: '*', element: <Navigate to="/" replace /> },
    ],
  },
]);

const App: React.FC = () => {
  return (
    // <React.StrictMode>
      <div className="App">
        <RouterProvider router={router} />
      </div>
    // </React.StrictMode>    
  );
};

export default App;