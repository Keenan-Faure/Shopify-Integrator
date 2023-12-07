import {useEffect} from 'react';
import Background from './Background';
import Setting_details from './semi-components/settings-details';
import '../CSS/page2.css';


function Page2(props)
{
    useEffect(()=> 
    {
        /* Ensure the model is shown */
        let navbar = document.getElementById("navbar");
        navbar.style.display = "block";

        /* Ensures the navbar is set correctly */
        let navigation = document.getElementById("navbar");
        window.onload = function(event)
        {
            navigation.style.left = "30%";
            navigation.style.position = "absolute";
            navigation.style.width = "70%";
            navigation.style.animation = "MoveLeft 1.2s ease";
        }

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

    }, []);

    return (
        <>
            <Background />
            <div className = "component1">
                <div className = "main-container">
                    <div className = "settings">
                        <div className = "app-settings">
                            <div className = "title">App Settings</div>
                            <div className = "setting">
                                <div className = "setting-title" style ={{top: '-6px'}}>Webhook Configuration
                                    <div className="info_icon" title="The forwarding url can be found in your ngrok dashboard."></div>
                                </div>
                                <div className = "setting-details description">Configures the Webhook required for the customers and order syncs to function correctly.</div>
                                <div className="webhook_div" action="/action_page.php" style= {{margin:  'auto',maxWidth: '300px'}}>
                                    <input type="text" placeholder = "Forwarding url..." name = "search2" />
                                    <button className = "button-on-off" type="submit">Create</button>
                                </div>
                            </div>

                            <div className = "setting" style = {{height: '240px', fontSize: '12px'}}>
                                <div className = "setting-title">Warehouse Location</div>
                                <div className = "setting-details description">Configures the location warehousing required for the products displayed</div>
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
                            
                            <Setting_details />
                            <Setting_details />
                        </div>
                        <div className = "shopify-settings">
                            <div className = "title">Shopify Settings</div>
                            <Setting_details />
                            <Setting_details />
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
                
            </div>
            
        </>
    );    

}    

export default Page2;