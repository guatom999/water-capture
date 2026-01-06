import { BrowserRouter, Routes, Route } from 'react-router-dom';
import LandingPage from './pages/LandingPage';
import SectionDetail from './components/SectionDetail';

function App() {
  return (

    <BrowserRouter>
      <Routes>
        <Route path="/" element={<LandingPage />} />
        <Route path="/markers/detail" element={<SectionDetail />} />
      </Routes>
    </BrowserRouter>
  )
}

export default App
