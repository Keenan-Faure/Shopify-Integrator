import {useEffect} from 'react';
import '../CSS/dashboard.css';
import image from '../media/icons8-shopify-50.png';
import Chart from "chart.js/auto";


function Dashboard()
{
    useEffect(()=> 
    {
        /* Ensures the navbar + model is set correctly */
        let navigation = document.getElementById("navbar");
        let logout = document.getElementById("logout");
        let header = document.querySelector('.header');
        let footer = document.querySelector('.footer');

        let ctx = document.getElementById('shopify_fetch_graph');
        let order_ctx = document.getElementById('shopify_order_graph');
        let graph1 = document.querySelector('.g1');
        let graph2 = document.querySelector('.g2');

        console.log(graph1);


        window.onload = function(event)
        {
            
            navigation.style.left = "0%";
            navigation.style.position = "absolute";
            navigation.style.width = "100%";
            navigation.style.display = "block";
            logout.style.display = "block"; 
        }

        /* Sets the initial Look of the Page */
        setTimeout(() => 
        {
            footer.style.display = "block";
            header.style.animation = "appear 1s ease-in";
            header.style.display = "block";

            
            graph1.style.animation = "appear 1s ease-in";
            graph1.style.display = "block";

            graph2.style.animation = "appear 1s ease-in";
            graph2.style.display = "block";
            
        }, 1200);


        const userName = localStorage.getItem('username');
        document.querySelector(".welcome_text").innerHTML = "Welcome back " + userName;

        /* logout */
        logout.addEventListener("click", () =>
        {
            logout.style.display = "none";
            navigation.style.display = "none";
            header.style.display = "none";
            footer.style.display = "none";

            /* Session Destroy */
            window.location.href = '/';
        });

        //Fetch Graph
        new Chart(ctx, {
            type: 'line',
            data: {
                labels: ['Monday', 'Tuesday', 'Wednesday', 'Thursday', 'Friday', 'Saturday', 'Sunday'],
                datasets: [{
                    label: '# of products fetched from Shopify',
                    data: [0, 20, 20, 10, 7, 0, 50],
                    borderWidth: 1,
                    borderColor: "rgba(255, 94, 0, 0.5)",
                    backgroundColor: "orange",
                    pointRadius: "5",
                }]
            },
            options: {
                scales: {
                    y: {
                        beginAtZero: true
                    }
                }
            }
        });

        //Order Graph
        new Chart(order_ctx, {
            type: 'bar',
            data: {
                labels: ['Monday', 'Tuesday', 'Wednesday', 'Thursday', 'Friday', 'Saturday', 'Sunday'],
                datasets: [{
                    label: 'Paid orders from Shopify (totals)',
                    data: [3400, 6100, 10000, 3456, 214, 0, 90],
                    borderWidth: 1,
                    borderColor: "rgba(85, 0, 255, 0.5);",
                    backgroundColor: "purple",
                    pointRadius: "5",
                    tension: 0.1,
                    fill: false
                },
                {
                    label: 'Unpaid orders from Shopify (totals)',
                    data: [800, 0, 2400, 4322, 100, 0, 40],
                    borderWidth: 1,
                    borderColor: "rgba(255,0,0,0.5);",
                    backgroundColor: "red",
                    pointRadius: "5",
                    tension: 0.1,
                    fill: false
                }]
            },
            options: {
                scales: {
                    y: {
                        beginAtZero: true
                    }
                }
            }
        });

    
        
        
        
    }, []);

    return (
        <div className = "dashboard" id = "dashboard">
            <div className = "container">

                <div className="header">
                    <div className="fetch_status_text">
                        Fetch Status
                        <div className="enabled_status">
                            <img className="logo" src= {image} />
                            <div className="logo_text">active</div>
                        </div>
                    </div>
                    <h2 className="welcome_text">Welcome back [Username]</h2>
                    <div className = "logout_button" id = "logout">Logout</div>
                </div>

                <div className="graph g1">
                    <canvas id="shopify_fetch_graph"></canvas>
                </div>

                <div className="graph g2">
                    <canvas id="shopify_order_graph"></canvas>
                </div>
                <div className="footer">
                    <p>Notifications</p>
                    <table>
                        <tbody>
                            <tr>
                                <td id="tr_logos">
                                    <div className="warning_logo"></div>
                                </td>
                                <td>Location-warehouse map needs to be configured</td>
                            </tr>
                            <tr>
                                <td id="tr_logos">
                                    <div className="warning_logo"></div>
                                </td>
                                <td>Queue needs to be enabled</td>
                            </tr>
                            <tr>
                                <td id="tr_logos">
                                    <div className="warning_logo"></div>
                                </td>
                                <td>Please visit settings</td>
                            </tr>
                        </tbody>
                    </table>
                </div>
            </div>
        </div>
                
    );
}


export default Dashboard;

/*

const ctx = document.getElementById('shopify_fetch_graph');
        const order_ctx = document.getElementById('shopify_order_graph');

        //Fetch Graph
        new Chart(ctx, {
            type: 'line',
            data: {
                labels: ['Monday', 'Tuesday', 'Wednesday', 'Thursday', 'Friday', 'Saturday', 'Sunday'],
                datasets: [{
                    label: '# of products fetched from Shopify',
                    data: [0, 20, 20, 10, 7, 0, 50],
                    borderWidth: 1,
                    borderColor: "rgba(255, 94, 0, 0.5)",
                    backgroundColor: "orange",
                    pointRadius: "5",
                }]
            },
            options: {
                scales: {
                    y: {
                        beginAtZero: true
                    }
                }
            }
        });

        //Order Graph
        new Chart(order_ctx, {
            type: 'bar',
            data: {
                labels: ['Monday', 'Tuesday', 'Wednesday', 'Thursday', 'Friday', 'Saturday', 'Sunday'],
                datasets: [{
                    label: 'Paid orders from Shopify (totals)',
                    data: [3400, 6100, 10000, 3456, 214, 0, 90],
                    borderWidth: 1,
                    borderColor: "rgba(85, 0, 255, 0.5);",
                    backgroundColor: "purple",
                    pointRadius: "5",
                    tension: 0.1,
                    fill: false
                },
                {
                    label: 'Unpaid orders from Shopify (totals)',
                    data: [800, 0, 2400, 4322, 100, 0, 40],
                    borderWidth: 1,
                    borderColor: "rgba(255,0,0,0.5);",
                    backgroundColor: "red",
                    pointRadius: "5",
                    tension: 0.1,
                    fill: false
                }]
            },
            options: {
                scales: {
                    y: {
                        beginAtZero: true
                    }
                }
            }
        });


*/

