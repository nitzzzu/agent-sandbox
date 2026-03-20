import ReactDOM from 'react-dom/client'
import App from './App'
import { applyTheme, getTheme } from './lib/theme/theme'
import './styles/index.css'

applyTheme(getTheme())

ReactDOM.createRoot(document.getElementById('root')!).render(<App />)
