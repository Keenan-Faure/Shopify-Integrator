import { Outlet, Link } from "react-router-dom";
import {useNavigate} from 'react-router-dom';
import '../CSS/navigation-bar.css';

/* Insert Logo below ↓ */
import MyImage from '../media/Screenshot.png';

const Layout = () => {

  const navigate = useNavigate();

  const navigateToProducts = () => 
  {
    navigate('/products');
  };
  const navigateToDashboard = () => 
  {
    navigate('/dashboard');
  };
  const navigateToOrders = () => 
  {
    navigate('/orders');
  };
  const navigateToCustomers = () => 
  {
    navigate('/customers');
  };
  const navigateToQueue = () => 
  {
    navigate('/queue');
  };
  const navigateToSettings = () => 
  {
    navigate('/settings');
  };

  return (
    <>
        <div className = "navbar" id = "navbar">

            <div className = "dropdown">
                <button onClick = {navigateToDashboard} className = "dropbtn">DashBoard</button>
            </div>

            <div className = "dropdown">
                <button onClick = {navigateToProducts} className = "dropbtn">Products</button>
                <div className = "dropdown-content">
                    <Link to="/importProducts">Import Product</Link>
                    <Link to="/exportProducts">Export Product</Link>
                    <Link to="/addProduct">Add Product</Link>
                </div>
            </div>

            <div className = "dropdown">
                <button onClick = {navigateToOrders} className = "dropbtn">Orders</button>
            </div>

            <div className = "dropdown">
                <button onClick = {navigateToCustomers} className = "dropbtn">Customers</button>
            </div>

            <div className = "dropdown">
                <button onClick = {navigateToQueue} className = "dropbtn">Queue</button>
            </div>

            <div className = "dropdown">
                <button onClick = {navigateToSettings} className = "dropbtn">Settings</button>
            </div>

            <div className = "avatar">
                <img className = "avatar" alt = "" src = {MyImage} />
            </div>


        </div>
      
        <Outlet />
    </>
  )
};
export default Layout