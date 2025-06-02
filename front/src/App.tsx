import { useState } from 'react'
import { BrowserRouter, Route, Routes } from 'react-router-dom'

import MainLayout from './layouts/MainLayout'
import LoginForm from './modules/LoginForm'
import RegisterForm from './modules/RegisterForm'
import GoodsPage from './pages/GoodsPage'
import Page404 from './pages/NotFoundPage'

function App() {
  const [admin, setAdmin] = useState(false)

  return (
    <BrowserRouter>
      <Routes>
        <Route path="/" element={<MainLayout admin={admin} />}>
          <Route path="lamps" element={<GoodsPage />} />
          <Route path="login" element={<LoginForm setAdmin={setAdmin} />} />
          <Route path="register" element={<RegisterForm />} />
          <Route path="*" element={<Page404 />} />
        </Route>
      </Routes>
    </BrowserRouter>
  )
}

export default App
