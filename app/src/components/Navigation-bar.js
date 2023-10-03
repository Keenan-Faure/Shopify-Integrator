// Imports Below 
import {BrowserRouter, Routes, Route} from "react-router-dom";
import {useEffect} from 'react';

/* Import links below ↓ */
import Layout from '../JS/Layout';
import Home from '../JS/DashBoard';
import NoPage from '../JS/Layout';
import Dashboard from "../JS/DashBoard";
import Products from '../JS/Products';
import Orders from '../JS/Orders';
import Settings from '../JS/Settings';
import Customers from '../JS/Customers';

// Import Style sheet below
import '../CSS/navigation-bar.css'
function Navigation_Bar(props)
{
    useEffect(()=> 
    {
        /* Change the style of navigation bar */
        let navigation = document.getElementById("navbar");
        let navbar = document.querySelectorAll(".dropbtn");
        let model = document.getElementById("model");

        /* The user clicks on 'Products' button */
        navbar[1].onclick = function(event)
        {
            if(navbar[1].onclick)
            {
                navigation.style.left = "30%";
                navigation.style.position = "absolute";
                navigation.style.width = "70%";
                navigation.style.animation = "MoveLeft 1.2s ease";
                model.style.animation = "none";
            } 
        }

        /* The user clicks on Buttons other than products or home */
        for(let i = 2; i < navbar.length; i++)
        {
            navbar[i].onclick = function(event)
            {
                if(navbar[i].onclick)
                {
                    navigation.style.left = "0%";
                    navigation.style.position = "relative";
                    navigation.style.width = "100%";
                    model.style.animation = "none";
                }
            }
        }

        /* The user clicks on 'Dashboard' button */
        navbar[0].onclick = function(event)
        {
            if(navbar[0].onclick)
            {
                navigation.style.animation = "MoveRight 1.2s ease";
                navigation.style.position = "fixed";
                navigation.style.left = "0%";
                navigation.style.width = "100%";
                model.style.animation = "none";
                model.style.display = "none";
            }
        }

    }, []);


    return (
        <div id = "navigation" style = {{display: props.Display}}>
            <BrowserRouter>
                <Routes>
                    <Route path = "/" element = {<Layout />}>
                        <Route index element = {<Home />}></Route>
                        <Route path = "dashboard" element = {<Dashboard />}></Route>
                        <Route path = "products" element = {<Products />}></Route>
                        <Route path = "orders" element = {<Orders />}></Route>
                        <Route path = "customers" element = {<Customers />}></Route>
                        <Route path = "settings" element = {<Settings />}></Route>

                        <Route path = "*" element = {<NoPage />}></Route>
                    </Route>
                </Routes>
            </BrowserRouter>
        </div>
    );
}

export default Navigation_Bar;
