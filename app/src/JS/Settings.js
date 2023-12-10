import {useEffect} from 'react';
import { createRoot } from 'react-dom/client';
import $ from 'jquery';
import Background from '../components/Background';
import Setting_details from '../components/semi-components/settings-details';

import '../CSS/page2.css';

function Settings()
{
    useEffect(()=> 
    {
        /* Ensures the page elements are set correctly */
        let navigation = document.getElementById("navbar");
        navigation.style.left = "25%";
        navigation.style.position = "absolute";
        navigation.style.width = "75%";


        /* Onclick for the Warehouse setting */
        let info_icon = document.querySelector(".info_icon");
        info_icon.addEventListener("click", () =>
        {
            alert("You will be forwarded to dashboard.ngrok.com");
            window.open('https://dashboard.ngrok.com', '_blank')
        })

        function fetchShopify() 
        {
            let shopify_locations = ["Japan", "Cape Town"];
            createLocationsDOM(shopify_locations);
            document.getElementById('fetch-button').disabled = true;
        }
    
        function createLocationsDOM(locations) 
        {
            let elements = document.querySelectorAll('.fill-able');
            for (let i = 0; i < elements.length; i++) 
            {
                let drop_down = document.createElement('select')
                let default_option = createOptions(true, "")
                drop_down.appendChild(default_option);
                for (let j = 0; j < locations.length; j++) 
                {
                    let option = createOptions(false, locations[j]);
                    drop_down.appendChild(option);
                }
                console.log(drop_down);
                elements[i].appendChild(drop_down);
            }
        }
    
        function createOptions(isDefault, location) 
        {
            let option = document.createElement('option');
            if (isDefault) 
            {
                option.setAttribute("value", location);
                option.innerHTML = "Select a location";
            } 
            else 
            {
                option.setAttribute("value", location);
                option.innerHTML = location;
            }
            return option;
        }

        /* Onclick for the Location setting */
        let fetch_button = document.getElementById("fetch-button");
        fetch_button.addEventListener("click", () => 
        {
            fetchShopify();
        });

        /*  API INITIAL-REQUEST for APP_SETTINGS*/
        const api_key = localStorage.getItem('api_key');
        $.ajaxSetup({ headers: { 'Authorization': 'ApiKey ' + api_key} });
        $.get("http://localhost:8080/api/settings", [], [])
        .done(function( _data) 
        {
            console.log(_data);
            
            let root;
            let _main = document.querySelector(".app-settings");
            let div = document.createElement("div");
            div.id = "a_settings";
            _main.appendChild(div);

            root = createRoot(div);
            root.render(_data.map((el, i) => <Setting_details key={`${el.title}_${i}`} Key={el.key} Description={el.description}
            Value={el.value} id={el.id}
            />))
            

            let setting_2 = document.getElementById("app_settings");
            for(let i = 0; i < _data.length; i++)
            {
                let div = document.createElement("button");
                div.className = "mini-setting";
                div.innerHTML = _data[i].key;
                setting_2.appendChild(div);
            }
            

            /* Scroll into View Button Event */
            let a_settings = document.getElementById("a_settings").childNodes;
            let app_button = document.getElementById("app_settings").children;
            console.log(app_button);
            console.log(a_settings);

            for(let i = 0; i < app_button.length; i++)
            {
                app_button[i].addEventListener("click", () =>
                {
                    a_settings[i].scrollIntoView({block: "center", behavior: 'smooth' });
                    setTimeout(() =>
                    {
                        a_settings[i].style.boxShadow = "0 0 20px rgb(173 216 230), 0 0 40px rgb(173 216 230), 0 0 60px rgb(173 216 230), 0 0 80px rgb(173 216 230), 0 0 80px rgb(173 216 230 / 10%)";
                        setTimeout(() =>
                        {
                            a_settings[i].style.boxShadow = "";
                        }, 1200)
                    }, 50);
                });
            }
            
        })
        .fail( function(xhr) { alert(xhr.responseText); });

        /*  API INITIAL-REQUEST for SHOPIFY_SETTINGS*/
        $.get("http://localhost:8080/api/shopify/settings", [], [])
        .done(function( _data) 
        {
            console.log(_data);
            
            let root;
            let _main = document.querySelector("._shopify");
            let div = document.createElement("div");
            div.id = "s_settings";
            _main.appendChild(div);

            root = createRoot(div);
            root.render(_data.map((el, i) => <Setting_details key={`${el.title}_${i}`} Key={el.key} Description={el.description}
            Value={el.value} id={el.id}
            />))

            let setting_2 = document.getElementById("shopify_settings");
            for(let i = 0; i < _data.length; i++)
            {
                let div = document.createElement("button");
                div.className = "mini-setting";
                div.innerHTML = _data[i].key;
                setting_2.appendChild(div);
            }

            /* Scroll into View Button Event */
            let s_settings = document.getElementById("s_settings").childNodes;
            let shop_button = document.getElementById("shopify_settings").children;
            console.log(shop_button);
            console.log(s_settings);
            for(let i = 0; i < shop_button.length; i++)
            {
                shop_button[i].addEventListener("click", () =>
                {
                    s_settings[i].scrollIntoView({block: "center", behavior: 'smooth' });
                    setTimeout(() =>
                    {
                        s_settings[i].style.boxShadow = "0 0 20px rgb(173 216 230), 0 0 40px rgb(173 216 230), 0 0 60px rgb(173 216 230), 0 0 80px rgb(173 216 230), 0 0 80px rgb(173 216 230 / 10%)";
                        setTimeout(() =>
                        {
                            s_settings[i].style.boxShadow = "";
                        }, 1200)
                    }, 50);
                    
                });
            }
            
        })
        .fail( function(xhr) { alert(xhr.responseText); });

        /* Submitted Setting Object */
        let edit = document.getElementById("edit");
        let confirm_line = document.querySelector(".confirm-line");
        edit.addEventListener("click", () =>
        {
            confirm_line.style.display = "block";
        });

        let confirm = document.getElementById("confirm");
        confirm.addEventListener("click", () =>
        {
            confirm_line.style.display = "none";
            let setting_main_title = document.querySelectorAll("._title");
            let setting_main_value = document.querySelectorAll("._input");

            let object = {};
            let _setting = {};
            for(let i = 0; i < setting_main_title.length; i++)
            {
                _setting = 
                {
                    key : setting_main_title[i].innerHTML,
                    value : setting_main_value[i].innerHTML
                }
                object[i] = _setting;
            }

            console.log(object);

            /*
            const api_key = localStorage.getItem('api_key');
            $.ajaxSetup({ headers: { 'Authorization': 'ApiKey ' + api_key} });
            $.post("http://localhost:8080/api/settings", JSON.stringify(object),[], 'json')
            .done(function( _data) 
            {
                console.log(_data);
            })
            .fail( function(xhr) 
            {
                alert(xhr.responseText);
            });
            */
        });
        

    }, []);

    return (
        <>
            <Background />
            <div className = "main-container">
                <div className = "settings">
                    <button className = "submiit" id = "edit" style = {{zIndex: '2', top: '55px'}}>Edit Settings</button>
                    <div className = "app-settings" style= {{position: 'relative', top:'15px'}}>
                        <div className = "title">App Settings</div>
                        <div className = "_app">
                            <div className = "setting">
                                <div className = "setting-title" style ={{top: '-6px'}}>Webhook Configuration
                                    <div className="info_icon" title="The forwarding url can be found in your ngrok dashboard."></div>
                                </div>
                                <div className = "setting-details description" style = {{textAlign: 'left', backgroundColor: 'transparent'}}>Configures the Webhook required for the customers and order syncs to function correctly.</div>
                                <div className="webhook_div" action="/action_page.php" style= {{margin:  'auto',maxWidth: '300px'}}>
                                    <input type="text" placeholder = "Forwarding url..." name = "search2" />
                                    <button className = "button-on-off" type="submit">Create</button>
                                </div>
                            </div>

                            <div className = "setting" style = {{height: '240px', fontSize: '12px'}}>
                                <div className = "setting-title">Warehouse Location</div>
                                <div className = "setting-details description" style = {{textAlign: 'left', backgroundColor: 'transparent'}}>Configures the location warehousing required for the products displayed</div>
                                <button className = "button-on-off" style = {{width: '90px'}}id="fetch-button">Fetch shopify locations</button>
                                <table style = {{left: '40%',top: '17px', marginBottom: '0px',fontSize: '13px'}}>
                                    <tbody>
                                        <tr>
                                            <th>Warehouse</th>
                                            <th>Shopify Location</th>
                                        </tr>
                                        <tr>
                                            <td>Cape Town Warehouse</td>
                                            <td className = "fill-able">
                                            </td>
                                        </tr>
                                        <tr>
                                            <td>Japan Warehouse</td>
                                            <td className = "fill-able">
                                            </td>
                                        </tr>
                                    </tbody> 
                                </table>
                            </div>
                        </div>
                        
                    </div>
                    <div className = "shopify-settings">
                        <div className = "title">Shopify Settings</div>
                        <div className = "_shopify"></div>
                    </div>   
                </div> 
                
            </div>
            <div className = "side-container">
                <div className = "settings-2">
                    <div className = "application"><i className = "a"/>Application Settings:</div>
                    <div id = "app_settings"></div>

                </div>
                <div className = "settings-2">
                    <div className = "application"><i className = "b"/>Spotify Settings:</div>
                    <div id = "shopify_settings"></div>
                </div>
            </div>
            <div className = "confirm-line">
                <button className="tablink" id = "confirm" style ={{left: '50%'/*, transform: 'translate(-50%)'*/}}>Save</button>
            </div>
            
        </>
    );
}

export default Settings;