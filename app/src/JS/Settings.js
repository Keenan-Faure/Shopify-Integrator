import React from 'react';
import {useEffect} from 'react';
import { createRoot } from 'react-dom/client';
import $ from 'jquery';
import Setting_details from '../components/semi-components/settings-details';
import Restriction_details from '../components/semi-components/restriction-details.js';
import Detailed_Restriction from '../components/semi-components/Settings/detailed-restrictions.js';
import Detailed_warehousing from '../components/semi-components/Settings/detailed_warehousing';
import Detailed_table from '../components/semi-components/Settings/detailed_table';

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
            alert("You will be forwarded to ngrok web-interface");
            window.open('http://localhost:8888/inspect/http', '_blank')
        })

        let _location_id = [];
        function fetchShopify() 
        {
            const api_key = localStorage.getItem('api_key');
            let shopify_locations  = [];
            let warehouse_locations  = [];

            /* Api-Request for shopify locations & warehouses */
            $.ajaxSetup({ headers: { 'Authorization': 'ApiKey ' + api_key} });
            $.get('http://localhost:8080/api/inventory/config', [], [])
            .done(function( _data) 
            {
                console.log(_data);

                for(let i = 0; i < _data.warehouses.length; i++)
                {
                    warehouse_locations[i] = _data.warehouses[i].name;
                }

                for(let i = 0; i < _data.shopify_locations.locations.length; i++)
                {
                    shopify_locations[i] = _data.shopify_locations.locations[i].name; 
                    _location_id[i] = _data.shopify_locations.locations[i].id; 
                }

                createLocationsDOM(shopify_locations);
                warehouseLocationsDOM(warehouse_locations);
                document.getElementById('fetch-button').disabled = true;

                /* changes the button */
                let button = document.getElementById("fetch-button");
                button.remove();

                let setting = document.querySelector(".setting");
                let new_button = document.createElement("button");
                new_button.className = "button-on-off";
                new_button.style.width = "90px";
                new_button.innerHTML = "Confirm Warehousing";
                setting.appendChild(new_button);

                new_button.addEventListener("click", () =>
                {
                    let object = {};
                    let select = document.querySelectorAll(".options");
                    let location_id = document.querySelector(".location-id");

                    object.location_id = location_id.innerHTML;
                    object.warehouse_name = select[0].options[select[0].selectedIndex].innerHTML;
                    object.shopify_warehouse_name = select[1].options[select[1].selectedIndex].innerHTML;

                    $.post("http://localhost:8080/api/inventory/map", JSON.stringify(object), [], 'json')
                    .done(function( _data) 
                    {
                        console.log(_data);
                        let info = document.getElementById("info-message");
                        info.innerHTML = "success";
                        info.style.display = "block";  
                        info.style.color = "white";
                        info.style.backgroundColor = "#1a5e12";
                        setTimeout(() => 
                        {
                            info.style.display = "none";  
                        }, 2000);
                    })
                    .fail( function(xhr) 
                    {
                        alert(xhr.responseText);
                    });
                });
                
            })
            .fail( function(xhr) { alert(xhr.responseText); });
        }
    
        function createLocationsDOM(locations) 
        {
            let elements = document.querySelectorAll('.fill-able');
            for (let i = 0; i < elements.length; i++) 
            {
                let drop_down = document.createElement('select');
                drop_down.className = "options 1";
                let default_option = createOptions(true, "")
                drop_down.appendChild(default_option);
                for (let j = 0; j < locations.length; j++) 
                {
                    let option = createOptions(false, locations[j]);
                    option.id = j;
                    drop_down.appendChild(option);
                }
                elements[i].appendChild(drop_down);
            }

            let select = document.querySelector(".options");

            select.addEventListener("change", () =>
            {
                let selectedOption = select.options[select.selectedIndex];
                let location_id = document.querySelector(".location-id");
                location_id.innerHTML = _location_id[parseInt(selectedOption.id)];
            });  
        }

        function warehouseLocationsDOM(locations) 
        {
            let elements = document.querySelectorAll('.warehouse');
            for (let i = 0; i < elements.length; i++) 
            {
                elements[i].innerHTML = "";
                let drop_down = document.createElement('select')
                drop_down.className = "options";
                let default_option = createOptions(true, "")
                drop_down.appendChild(default_option);
                for (let j = 0; j < locations.length; j++) 
                {
                    let option = createOptions(false, locations[j]);
                    drop_down.appendChild(option);
                }
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
            if(document.getElementById("a_settings") != null && document.getElementById("app_settings") != null)
            {
                // Remove old elements
                document.getElementById("a_settings").remove();
                document.getElementById("app_settings").remove();

                let settings_2 = document.getElementById("a-s");
                let div = document.createElement("div");
                div.id = "a_settings";
                _main.appendChild(div);
                root = createRoot(div);
                root.render(_data.map((el, i) => <Setting_details key={`${el.title}_${i}`} Key={el.field_name} Description={el.description}
                Value={el.value} id={el.id} Title = {el.key}
                />))

                //create new setting_2
                let setting_2 = document.createElement("div");
                setting_2.id = "app_settings";
                settings_2.appendChild(setting_2);

                for(let i = 0; i < _data.length; i++)
                {
                    let div = document.createElement("button");
                    div.className = "mini-setting";
                    div.innerHTML = _data[i].field_name;
                    setting_2.appendChild(div);
                }
            
                /* Scroll into View Button Event */
                let a_settings = div.childNodes;
                let app_button = setting_2.children;
                for(let i = 0; i < app_button.length; i++)
                {
                    app_button[i].addEventListener("click", () =>
                    {
                        a_settings[i].scrollIntoView({block: "center", behavior: 'smooth' });
                        setTimeout(() =>
                        {
                            a_settings[i].style.boxShadow = "0 0 20px rgb(173 216 230), 0 0 40px rgb(173 216 230), 0 0 60px rgb(173 216 230), 0 0 80px rgb(173 216 230), 0 0 80px rgb(173 216 230 / 10%)";
                            setTimeout(() => { a_settings[i].style.boxShadow = ""; }, 1200)
                        }, 50);
                    });
                }
            }
            else 
            {
                let div = document.createElement("div");
                div.id = "a_settings";
                _main.appendChild(div);
                root = createRoot(div);
                root.render(_data.map((el, i) => <Setting_details key={`${el.title}_${i}`} Key={el.field_name} Description={el.description}
                Value={el.value} id={el.id} Title = {el.key}
                />))

                let setting_2 = document.getElementById("app_settings");
                for(let i = 0; i < _data.length; i++)
                {
                    let div = document.createElement("button");
                    div.className = "mini-setting";
                    div.innerHTML = _data[i].field_name;
                    setting_2.appendChild(div);
                }
            
                /* Scroll into View Button Event */
                let a_settings = div.childNodes;
                let app_button = setting_2.children;

                for(let i = 0; i < app_button.length; i++)
                {
                    app_button[i].addEventListener("click", () =>
                    {
                        a_settings[i].scrollIntoView({block: "center", behavior: 'smooth' });
                        setTimeout(() =>
                        {
                            a_settings[i].style.boxShadow = "0 0 20px rgb(173 216 230), 0 0 40px rgb(173 216 230), 0 0 60px rgb(173 216 230), 0 0 80px rgb(173 216 230), 0 0 80px rgb(173 216 230 / 10%)";
                            setTimeout(() => { a_settings[i].style.boxShadow = ""; }, 1200)
                        }, 50);
                    })
                }
            }
        
        })
        .fail( function(xhr) 
        { 
            alert(xhr.responseText); 
        });

        /*  API INITIAL-REQUEST for SHOPIFY_SETTINGS*/
        $.get("http://localhost:8080/api/shopify/settings", [], [])
        .done(function( _data) 
        {
            console.log(_data);
            
            let root;
            let _main = document.querySelector("._shopify");
            if(document.getElementById("s_settings") != null && document.getElementById("shopify_settings") != null)
            {
                document.getElementById("s_settings").remove();
                document.getElementById("shopify_settings").remove();

                let settings_2 = document.getElementById("s-s");
                let div = document.createElement("div");
                div.id = "s_settings";
                _main.appendChild(div);

                root = createRoot(div);
                root.render(_data.map((el, i) => <Setting_details key={`${el.title}_${i}`} Key={el.field_name} Description={el.description}
                Value={el.value} id={el.id} Title = {el.key}
                />))

                //create new setting_2
                let setting_2 = document.createElement("div");
                setting_2.id = "shopify_settings";
                settings_2.appendChild(setting_2);

                for(let i = 0; i < _data.length; i++)
                {
                    let div = document.createElement("button");
                    div.className = "mini-setting";
                    div.innerHTML = _data[i].field_name;
                    setting_2.appendChild(div);
                }

                /* Scroll into View Button Event */
                let s_settings = div.childNodes;
                let shop_button = setting_2.children;
                for(let i = 0; i < shop_button.length; i++)
                {
                    shop_button[i].addEventListener("click", () =>
                    {
                        s_settings[i].scrollIntoView({block: "center", behavior: 'smooth' });
                        setTimeout(() =>
                        {
                            s_settings[i].style.boxShadow = "0 0 20px rgb(173 216 230), 0 0 40px rgb(173 216 230), 0 0 60px rgb(173 216 230), 0 0 80px rgb(173 216 230), 0 0 80px rgb(173 216 230 / 10%)";
                            setTimeout(() => { s_settings[i].style.boxShadow = ""; }, 1200)
                        }, 50);
                    });
                }
                
            }
            else 
            {
                let div = document.createElement("div");
                div.id = "s_settings";
                _main.appendChild(div);

                root = createRoot(div);
                root.render(_data.map((el, i) => <Setting_details key={`${el.title}_${i}`} Key={el.field_name} Description={el.description}
                Value={el.value} id={el.id} Title = {el.key}
                />))

                let setting_2 = document.getElementById("shopify_settings");
                for(let i = 0; i < _data.length; i++)
                {
                    let div = document.createElement("button");
                    div.className = "mini-setting";
                    div.innerHTML = _data[i].field_name;
                    setting_2.appendChild(div);
                }

                /* Scroll into View Button Event */
                let s_settings = div.childNodes;
                let shop_button = setting_2.children;
                for(let i = 0; i < shop_button.length; i++)
                {
                    shop_button[i].addEventListener("click", () =>
                    {
                        s_settings[i].scrollIntoView({block: "center", behavior: 'smooth' });
                        setTimeout(() =>
                        {
                            s_settings[i].style.boxShadow = "0 0 20px rgb(173 216 230), 0 0 40px rgb(173 216 230), 0 0 60px rgb(173 216 230), 0 0 80px rgb(173 216 230), 0 0 80px rgb(173 216 230 / 10%)";
                            setTimeout(() => { s_settings[i].style.boxShadow = ""; }, 1200)
                        }, 50);
                        
                    })
                }
            }
        })
        .fail( function(xhr) 
        { 
            alert(xhr.responseText); 
        });

        /* Restrictions Settings */
        $.ajaxSetup({ headers: { 'Authorization': 'ApiKey ' + api_key} });
        $.get("http://localhost:8080/api/fetch/restriction", [], [], 'json')
        .done(function( _data) 
        {
            
            let root;
            let _main = document.querySelector(".restriction-settings");

            if(document.getElementById("r_settings") != null && document.getElementById("restriction_settings") != null)
            {
                document.getElementById("r_settings").remove();
                document.getElementById("restriction_settings").remove();

                let settings_2 = document.getElementById("r-s");
                let div = document.createElement("div");
                div.id = "r_settings";
                _main.appendChild(div);
                root = createRoot(div);
                root.render(<Restriction_details key={`${_data.id}`} restriction = { _data.map((el, i) => <Detailed_Restriction key={`${el.title}_${i}`} 
                Name = {el.field[0].toUpperCase() + el.field.slice(1)} Flag = {el.flag}
                />)}/>)

                //create new setting_2
                let setting_2 = document.createElement("div");
                setting_2.id = "restriction_settings";
                settings_2.appendChild(setting_2);

                for(let i = 0; i < 1; i++)
                {
                    let div = document.createElement("button");
                    div.className = "mini-setting";
                    div.innerHTML = "Fetch Restrictions";
                    setting_2.appendChild(div);
                }
            
                /* Scroll into View Button Event */
                let r_settings = div.childNodes;
                let rest_button = setting_2.children;
                for(let i = 0; i < 1; i++)
                {
                    rest_button[i].addEventListener("click", () =>
                    {
                        r_settings[i].scrollIntoView({block: "center", behavior: 'smooth' });
                        setTimeout(() =>
                        {
                            r_settings[i].style.boxShadow = "0 0 20px rgb(173 216 230), 0 0 40px rgb(173 216 230), 0 0 60px rgb(173 216 230), 0 0 80px rgb(173 216 230), 0 0 80px rgb(173 216 230 / 10%)";
                            setTimeout(() => { r_settings[i].style.boxShadow = ""; }, 1200)
                        }, 50);
                    });
                }
            }
            else 
            {
                let settings_2 = document.getElementById("r-s");
                let div = document.createElement("div");
                div.id = "r_settings";
                _main.appendChild(div);
                root = createRoot(div);
                root.render(<Restriction_details key={`${_data.id}`} restriction = { _data.map((el, i) => <Detailed_Restriction key={`${el.title}_${i}`} 
                Name = {el.field[0].toUpperCase() + el.field.slice(1)} Flag = {el.flag}
                />)}/>)

                //create new setting_2
                let setting_2 = document.createElement("div");
                setting_2.id = "restriction_settings";
                settings_2.appendChild(setting_2);

                for(let i = 0; i < 1; i++)
                {
                    let div = document.createElement("button");
                    div.className = "mini-setting";
                    div.innerHTML = "Fetch Restrictions";
                    setting_2.appendChild(div);
                }
            
                /* Scroll into View Button Event */
                let r_settings = div.childNodes;
                let rest_button = setting_2.children;
                for(let i = 0; i < 1; i++)
                {
                    rest_button[i].addEventListener("click", () =>
                    {
                        r_settings[i].scrollIntoView({block: "center", behavior: 'smooth' });
                        setTimeout(() =>
                        {
                            r_settings[i].style.boxShadow = "0 0 20px rgb(173 216 230), 0 0 40px rgb(173 216 230), 0 0 60px rgb(173 216 230), 0 0 80px rgb(173 216 230), 0 0 80px rgb(173 216 230 / 10%)";
                            setTimeout(() => { r_settings[i].style.boxShadow = ""; }, 1200)
                        }, 50);
                    });
                }
            } 
        })
        .fail( function(xhr) { alert(xhr.responseText); });

        /* Restrictions Settings */
        $.ajaxSetup({ headers: { 'Authorization': 'ApiKey ' + api_key} });
        $.get("http://localhost:8080/api/push/restriction", [], [], 'json')
        .done(function( _data) 
        {

            let root;
            let _main = document.querySelector(".restriction-settings-2");

            if(document.getElementById("r_settings2") != null && document.getElementById("restriction_settings2") != null) 
            {
                document.getElementById("r_settings2").remove();
                document.getElementById("restriction_settings2").remove();

                let settings_2 = document.getElementById("r-s");
                let div = document.createElement("div");
                div.id = "r_settings2";
                _main.appendChild(div);
                root = createRoot(div);
                root.render(<Restriction_details key={`${_data.id}`} restriction = { _data.map((el, i) => <Detailed_Restriction key={`${el.title}_${i}`} 
                Name = {el.field[0].toUpperCase() + el.field.slice(1)} Flag = {el.flag}
                />)}/>)

                //create new setting_2
                let setting_2 = document.createElement("div");
                setting_2.id = "restriction_settings2";
                settings_2.appendChild(setting_2);

                for(let i = 0; i < 1; i++)
                {
                    let div = document.createElement("button");
                    div.className = "mini-setting";
                    div.innerHTML = "Push Restrictions";
                    setting_2.appendChild(div);
                }
            
                /* Scroll into View Button Event */
                let r_settings = div.childNodes;
                let rest_button = setting_2.children;
                for(let i = 0; i < 1; i++)
                {
                    rest_button[i].addEventListener("click", () =>
                    {
                        r_settings[i].scrollIntoView({block: "center", behavior: 'smooth' });
                        setTimeout(() =>
                        {
                            r_settings[i].style.boxShadow = "0 0 20px rgb(173 216 230), 0 0 40px rgb(173 216 230), 0 0 60px rgb(173 216 230), 0 0 80px rgb(173 216 230), 0 0 80px rgb(173 216 230 / 10%)";
                            setTimeout(() => { r_settings[i].style.boxShadow = ""; }, 1200)
                        }, 50);
                    });
                }
            }
            else 
            {
                let settings_2 = document.getElementById("r-s");
                let div = document.createElement("div");
                div.id = "r_settings2";
                _main.appendChild(div);
                root = createRoot(div);
                root.render(<Restriction_details key={`${_data.id}`} restriction = { _data.map((el, i) => <Detailed_Restriction key={`${el.title}_${i}`} 
                Name = {el.field[0].toUpperCase() + el.field.slice(1)} Flag = {el.flag}
                />)}/>)

                //create new setting_2
                let setting_2 = document.createElement("div");
                setting_2.id = "restriction_settings2";
                settings_2.appendChild(setting_2);

                for(let i = 0; i < 1; i++)
                {
                    let div = document.createElement("button");
                    div.className = "mini-setting";
                    div.innerHTML = "Push Restrictions";
                    setting_2.appendChild(div);
                }
            
                /* Scroll into View Button Event */
                let r_settings = div.childNodes;
                let rest_button = setting_2.children;
                for(let i = 0; i < 1; i++)
                {
                    rest_button[i].addEventListener("click", () =>
                    {
                        r_settings[i].scrollIntoView({block: "center", behavior: 'smooth' });
                        setTimeout(() =>
                        {
                            r_settings[i].style.boxShadow = "0 0 20px rgb(173 216 230), 0 0 40px rgb(173 216 230), 0 0 60px rgb(173 216 230), 0 0 80px rgb(173 216 230), 0 0 80px rgb(173 216 230 / 10%)";
                            setTimeout(() => { r_settings[i].style.boxShadow = ""; }, 1200)
                        }, 50);
                    });
                }
            }
            
        })
        .fail( function(xhr) { alert(xhr.responseText); });

        let confirm_line = document.querySelector(".confirm-line");
    
        let confirm = document.getElementById("confirm");
        confirm.addEventListener("click", () =>
        {
            confirm_line.style.display = "none";
            let app_div = document.querySelector(".app-settings");
            let setting_main_title = app_div.querySelectorAll("._title"); let setting_main_value = app_div.querySelectorAll("._input");
            let setting_default_value = app_div.querySelectorAll("._value");

            let shop_div = document.querySelector(".shopify-settings");
            let setting_main_title2 = shop_div.querySelectorAll("._title"); let setting_main_value2 = shop_div.querySelectorAll("._input");
            let setting_default_value2 = shop_div.querySelectorAll("._value");

            let button_true = document.querySelectorAll(".true");
            let button_false = document.querySelectorAll(".false");

            let app_object = [];
            let _shopify_object = [];
            let _setting = {};
            for(let i = 0; i < setting_main_title.length ; i++)
            {
                if(setting_main_value[i].style.display == "block")
                {
                    //IF setting is empty, use the original value stored in default
                    if(setting_main_value[i].value == "")
                    {
                        _setting = 
                        {
                            key : setting_main_title[i].innerHTML,
                            value : setting_default_value[i].innerHTML
                        }
                        app_object[i] = _setting;
                    }
                    //use the new version
                    else 
                    {
                        _setting = 
                        {
                            key : setting_main_title[i].innerHTML,
                            value : setting_main_value[i].value
                        }
                        app_object[i] = _setting;
                    }
                }
                else 
                {
                    // IF neither buttons are selected, use the default value
                    if(button_false[i].className == "false" && button_true[i].className == "true")
                    {
                        _setting = 
                        {
                            key : setting_main_title[i].innerHTML,
                            value : setting_default_value[i].innerHTML
                        }
                        app_object[i] = _setting;
                    }   
                    //if button false is selected, use false
                    else if(button_false[i].className != "false")
                    {
                        _setting = 
                        {
                            key : setting_main_title[i].innerHTML,
                            value : button_false[i].innerHTML.toLowerCase()
                        }
                        app_object[i] = _setting;
                    }
                    //if button true is slelected, use true
                    else if(button_true[i].className != "true")
                    {
                        _setting = 
                        {
                            key : setting_main_title[i].innerHTML,
                            value : button_true[i].innerHTML.toLowerCase()
                        }
                        app_object[i] = _setting;
                    }
                }   
            }

            let count = 0;
            for(let i = 0 ; i < setting_main_title2.length; i++)
            {
                if(setting_main_value2[count].style.display == "block")
                {
                    //IF setting is empty, use the original value stored in default
                    if(setting_main_value2[count].value == "")
                    {
                        _setting = 
                        {
                            key : setting_main_title2[i].innerHTML,
                            value : setting_default_value2[i].innerHTML
                        }
                        _shopify_object[count] = _setting;
                    }
                    //use the new version
                    else 
                    {
                        _setting = 
                        {
                            key : setting_main_title2[i].innerHTML,
                            value : setting_main_value2[i].value
                        }
                        _shopify_object[count] = _setting;
                    }
                }
                else 
                {
                    // IF neither buttons are selected, use the default value
                    if(button_false[count].className == "false" && button_true[i].className == "true")
                    {
                        _setting = 
                        {
                            key : setting_main_title2[i].innerHTML,
                            value : setting_default_value2[i].innerHTML
                        }
                        _shopify_object[count] = _setting;
                    }   
                    //if button false is selected, use false
                    else if(button_false[count].className != "false")
                    {
                        _setting = 
                        {
                            key : setting_main_title2[i].innerHTML,
                            value : button_false[i].innerHTML.toLowerCase()
                        }
                        _shopify_object[count] = _setting;
                    }
                    //if button true is slelected, use true
                    else if(button_true[count].className != "true")
                    {
                        _setting = 
                        {
                            key : setting_main_title2[i].innerHTML,
                            value : button_true[i].innerHTML.toLowerCase()
                        }
                        _shopify_object[count] = _setting;
                    }
                }
                count += 1;
            }

            /* Restrictions */
            let fetch_restrict = [];
            let push_restrict = [];

            let rest1 = document.getElementById("r_settings"); let rest2 = document.getElementById("r_settings2");
            let fetch_field = rest1.querySelectorAll(".restriction_name");
            let switch_fetch = rest1.querySelectorAll ("#lide"); let switch_push = rest2.querySelectorAll("#lide");
            let push_field = rest2.querySelectorAll(".restriction_name");

            let _fetch = {};
            let _push = {};
            let _c = 0;
            let c_ = 0;
            for(let i = 0; i < fetch_field.length; i++)
            {
                if(switch_fetch[i].checked == true)
                {
                    _fetch = 
                    {
                        field : fetch_field[i].innerHTML.toLowerCase(),
                        flag : "app"
                    }
                    fetch_restrict[_c] = _fetch;
                }
                else 
                {
                    _fetch = 
                    {
                        field : fetch_field[i].innerHTML.toLowerCase(),
                        flag : "shopify"
                    }
                    fetch_restrict[_c] = _fetch;
                }
                _c += 1;
            }

            for(let i = 0; i < push_field.length; i++)
            {
                if(switch_push[i].checked == true)
                {
                    _push = 
                    {
                        field : fetch_field[i].innerHTML.toLowerCase(),
                        flag : "app"
                    }
                    push_restrict[c_] = _push;
                }
                else 
                {
                    _push = 
                    {
                        field : fetch_field[i].innerHTML.toLowerCase(),
                        flag : "shopify"
                    }
                    push_restrict[c_] = _push;
                }
                c_ += 1;
            }

            const api_key = localStorage.getItem('api_key');

            console.log(fetch_restrict);
            console.log(push_restrict);

            /* Fetch Restriction Put Setting */
            $.ajaxSetup({ headers: { 'Authorization': 'ApiKey ' + api_key}});
            $.ajax({ type: 'PUT', url: "http://localhost:8080/api/fetch/restriction", contentType: 'json', data: JSON.stringify(fetch_restrict)})
            .done(function (_data) 
            {
                console.log(_data);
                let info = document.getElementById("info-message");
                info.innerHTML = "success";
                info.style.display = "block";  
                info.style.color = "white";
                info.style.backgroundColor = "#1a5e12";
                setTimeout(() => 
                {
                    info.style.display = "none";  
                }, 2000);
            })
            .fail( function(xhr) 
            {
                alert(xhr.responseText);
            });

            /* Push Restriction Put Setting */
            $.ajax({ type: 'PUT', url: "http://localhost:8080/api/push/restriction", contentType: 'json', data: JSON.stringify(push_restrict)})
            .done(function (_data) 
            {
                console.log(_data);
                let info = document.getElementById("info-message");
                info.innerHTML = "success";
                info.style.display = "block";  
                info.style.color = "white";
                info.style.backgroundColor = "#1a5e12";
                setTimeout(() => 
                {
                    info.style.display = "none";  
                }, 2000);
            })
            .fail( function(xhr) 
            {
                alert(xhr.responseText);
            });
            
            
            /* App Setting Put Request */
            $.ajax({ type: 'PUT', url: "http://localhost:8080/api/settings", contentType: 'json', data: JSON.stringify(app_object)})
            .done(function (_data) 
            {
                console.log(_data);
                let info = document.getElementById("info-message");
                info.innerHTML = "success";
                info.style.display = "block";  
                info.style.color = "white";
                info.style.backgroundColor = "#1a5e12";
                setTimeout(() => 
                {
                    info.style.display = "none";  
                }, 2000);
            })
            .fail( function(xhr) 
            {
                alert(xhr.responseText);
            });

            /* Shopify Setting Put Request */
            $.ajax({ type: 'PUT', url: "http://localhost:8080/api/shopify/settings",
            contentType: 'json', data: JSON.stringify(_shopify_object)})
            .done(function (_data) 
            {
                console.log(_data);
                let info = document.getElementById("info-message");
                info.innerHTML = "success";
                info.style.display = "block";  
                info.style.color = "white";
                info.style.backgroundColor = "#1a5e12";
                setTimeout(() => 
                {
                    info.style.display = "none";  
                }, 2000);
            })
            .fail( function(xhr) 
            {
                alert(xhr.responseText);
            });
        });

        /* Webhook Setting */
        let button = document.getElementById("_webhook");
        button.addEventListener("click", () =>
        {
            button.remove();
            let webhook_button = document.createElement("button");
            webhook_button.className = "webhook-button";
            webhook_button.innerHTML = "Copy Webhook";

            let domain_name = document.getElementById("domain_name").value;

            if(domain_name == "")
            {
                alert("Domain name field is empty...")
            }
            else 
            {
                let setting = document.getElementById("setting1");
                setting.appendChild(webhook_button);
                let domain =  
                {
                    domain: domain_name
                }
                
                let copyText;
                let data;

                const api_key = localStorage.getItem('api_key');
                $.ajaxSetup({ headers: { 'Authorization': 'ApiKey ' + api_key} });
                
                $.post("http://localhost:8080/api/settings/webhook", JSON.stringify(domain), [], 'json')
                .done(function( _data) 
                {
                    console.log(_data);

                    let info = document.getElementById("info-message");
                    info.innerHTML = "success";
                    info.style.display = "block"; 
                    info.style.color = "white";
                    info.style.backgroundColor = "#1a5e12";
                    setTimeout(() => 
                    {
                        info.style.display = "none";  
                    }, 2000);

                    data = _data.message;
                })
                .fail( function(xhr) 
                {
                    alert(xhr.responseText);
                });
                
                webhook_button.addEventListener("click", () =>
                { 
                    copyText = data;
                    navigator.clipboard.writeText(copyText);
                    webhook_button.innerHTML = "Copied!";
                });
            } 
        });

        /* Adds a new Warehouse */
        let warehouse_button = document.getElementById("new_warehouse");
        warehouse_button.addEventListener("click", () =>
        {
            let _input = document.getElementById("ware_").value;

            let body =
            {
                name: _input
            }

            const api_key = localStorage.getItem('api_key');
            $.ajaxSetup({ headers: { 'Authorization': 'ApiKey ' + api_key} });
            $.post("http://localhost:8080/api/inventory/warehouse?reindex=false", JSON.stringify(body), [], 'json')
            .done(function( _data) 
            {
                console.log(_data);
                let info = document.getElementById("info-message");
                info.innerHTML = _data.message;
                info.style.display = "block";  
                info.style.color = "white";
                info.style.backgroundColor = "#1a5e12"; 
                setTimeout(() => 
                {
                    info.style.display = "none";  
                }, 2000);           
            })
            .fail( function(xhr) 
            {
                alert(xhr.responseText);
            });
        });

        setTimeout(() =>
        {
            let setting_main = document.querySelectorAll(".setting_main");

            for(let i=0; i<setting_main.length; i++)
            {
                let setting_default_value = setting_main[i].querySelector("._value");
                let setting_input = setting_main[i].querySelector("._input");
                let true_false = setting_main[i].querySelector(".true-false");
                if(setting_default_value.innerHTML == "true" || setting_default_value.innerHTML == "false")
                {
                    setting_input.style.display = "none";
                }
                else
                {
                    true_false.style.display = "none";
                }
            }
        }, 250);

        //If the user changes any of the input fields/buttons show the save button
        setTimeout(() =>
        {
            let setting_main_value = document.querySelectorAll("._input");
            let button_true = document.querySelectorAll(".true");
            let button_false = document.querySelectorAll(".false");

            //If button
            for(let i = 0; i < button_true.length; i++)
            {
                button_true[i].addEventListener("click", () => 
                {
                    confirm_line.style.display = "block";
                });

                button_false[i].addEventListener("click", () => 
                {
                    confirm_line.style.display = "block";
                });
            }

            //If input field 
            for(let i = 0; i < setting_main_value.length; i++)
            {
                setting_main_value[i].addEventListener("input", () => 
                {
                    confirm_line.style.display = "block";
                });
            }
        }, 250);


        /* Displays the warehouse map */
        let warehouse_map = document.getElementById("warehouse_map");
        let side = document.querySelector(".side-container");
        let main = document.querySelector(".main-container");
        let return_button = document.querySelector(".rtn-button");
        let _main = document.getElementById("_main");

        warehouse_map.addEventListener("click", () =>
        {
            let map_container = document.querySelector(".warehouse-mapp");
            return_button.style.display = "block";
            map_container.style.display = "block";
            side.style.display = "none";
            main.style.display = "none";
            navigation.style.display = "none";

            $.get("http://localhost:8080/api/inventory/map", [], [], 'json')
            .done(function( _data) 
            {
                console.log(_data);
                let info = document.getElementById("info-message");
                info.innerHTML = _data;
                info.style.display = "block";
                info.style.color = "white";
                setTimeout(() => 
                {
                    info.style.display = "none";  
                }, 2000);
                if(document.querySelector(".warehouse-mapp") != null)
                {
                    document.querySelector(".warehouse-mapp").remove();

                    let root;
                    let div = document.createElement("div");
                    div.className = "warehouse-mapp";
                    div.style.display = "block";
                    _main.appendChild(div);
                    root = createRoot(div);
                    root.render(<Detailed_table key={_data.title} table={_data.map((el, i) => <Detailed_warehousing Created_At={el.created_at} 
                    Shopify_Location_ID={el.shopify_location_id} id={el.id} 
                    Shopify_Warehouse_Name ={el.shopify_warehouse_name} Warehouse_Name={el.warehouse_name} />)}
                    />)
                }
                else 
                {
                    let root;
                    let div = document.createElement("div");
                    div.className = "warehouse-mapp";
                    div.style.display = "block";
                    _main.appendChild(div);
                    root = createRoot(div);
                    root.render(<Detailed_table key={_data.title} table={_data.map((el, i) => <Detailed_warehousing Created_At={el.created_at} 
                    Shopify_Location_ID={el.shopify_location_id} id={el.id} 
                    Shopify_Warehouse_Name ={el.shopify_warehouse_name} Warehouse_Name={el.warehouse_name} />)}
                    />)
                }
            })
            .fail( function(xhr) 
            {
                alert(xhr.responseText);
            });
        });

        return_button.addEventListener("click", () =>
        {
            let map_container = document.querySelector(".warehouse-mapp");
            let navigation = document.getElementById("navbar");

            navigation.style.display = "block";
            return_button.style.display = "none";
            map_container.style.display = "none";
            setTimeout(() => 
            {
                side.style.display = "block";
                main.style.display = "block";
            }, 800);
        });

        setTimeout(() => 
        {
            /* Switch event Change*/
            let _switch = document.querySelectorAll(".switch");
            let confirm_line = document.querySelector(".confirm-line");
            for(let i = 0; i < _switch.length; i++)
            {
                _switch[i].addEventListener("click", () =>
                {
                    confirm_line.style.display = "block";
                });
            }
        }, 300);
        

    }, []);

    return (
        <div id = "_main">
            <div className = "main-container">
                <div className = "settings">
                    <div className = "app-settings" style= {{position: 'relative', top:'15px'}}>
                        <div className = "title" title = "Application Settings" style = {{marginTop: '20px'}}>App Settings</div>
                        <div className = "_app">
                            <div className = "setting" style = {{height: '220px', fontSize: '12px'}}>
                                <div className = "setting-title">Warehouse Location</div>
                                <div className = "setting-details description" style = {{textAlign: 'center', backgroundColor: 'transparent'}}>Configures the location warehousing required for the products displayed</div>
                                <button className = "mini-setting" style ={{width: '20%'}} id = "warehouse_map">Current Warehouse Map</button>
                                <button className = "button-on-off" style = {{width: '90px'}}id="fetch-button">Fetch shopify locations</button>
                                
                                <table style = {{left: '40%',top: '17px', marginBottom: '0px',fontSize: '13px'}}>
                                    <tbody>
                                        <tr>
                                            <th>Warehouse</th>
                                            <th>Shopify Location</th>
                                            <th>Location ID</th>
                                        </tr>
                                        <tr>
                                            <td className = "warehouse">Warehouse</td>
                                            <td className = "fill-able"></td>
                                            <td className = "location-id"></td>
                                        </tr>
                                    </tbody>
                                </table>
                            </div>
                            
                            <div className = "setting_ware">
                                <div className = "setting-title" style ={{top: '-14px'}}>Add Warehouse</div>
                                <div className = "setting-details description" style = {{textAlign: 'center', backgroundColor: 'transparent'}}>
                                    Adds a new Warehouse to the list of existing Warehouses
                                </div>
                                <div className="webhook_div" style= {{margin:  'auto',maxWidth: '300px'}}>
                                    <input id = "ware_" type="text" placeholder = "Warehouse Name..." name = "new_warehouse" style= {{color: 'black'}}/>
                                    <button id = "new_warehouse" className = "button-on-off" type="button">Create</button>
                                </div>
                            </div>

                            <div className = "setting" id = "setting1">
                                <div className = "setting-title" style ={{top: '-14px'}}>Webhook Configuration
                                    <div className="info_icon" title="The forwarding url can be found in your ngrok dashboard."></div>
                                </div>
                                <div className = "setting-details description" style = {{textAlign: 'left', backgroundColor: 'transparent'}}>Configures the Webhook required for the customers and order syncs to function correctly.</div>
                                <div className="webhook_div" style= {{margin:  'auto',maxWidth: '300px'}}>
                                    <input type="text" placeholder = "Domain Name..." id = "domain_name" name = "domain_name" style= {{color: 'black'}}/>
                                    <button id = "_webhook" className = "button-on-off" type="button">Create</button>
                                </div>
                            </div>
                        </div>
                        
                    </div>
                    <div className = "shopify-settings">
                        <div className = "title" title = "Shopify Settings">Shopify Settings</div>
                        <div className = "_shopify"></div>
                    </div> 
                    <div className = "restriction-settings">
                        <div className = "title" title = "Restricts certain information from Shopify Fetch">Fetch Restrictions</div>
                    </div> 

                    <div className = "restriction-settings-2">
                        <div className = "title" title = "Restricts certain information when pushing to Shopify">Push Restrictions</div>
                    </div>

                </div> 
                
            </div>
            <div className = "side-container">
                <div className = "settings-2" id = "a-s">
                    <div className = "application"><i className = "a"/>Application Settings:</div>
                    <div id = "app_settings"></div>

                </div>
                <div className = "settings-2" id = "s-s" style ={{paddingTop: '0px'}}>
                    <div className = "application"><i className = "b"/>Spotify Settings:</div>
                    <div id = "shopify_settings"></div>
                </div>
                <div className = "settings-2" id = "r-s" style ={{paddingTop: '0px'}}>
                    <div className = "application"><i className = "c"/>Restriction Settings:</div>
                    <div id = "restrictions"></div>
                </div>
            </div>
            <div className = "confirm-line">
                <button className="tablink" id = "confirm" style ={{left: '50%'}}>Save</button>
            </div>
            <div className = 'rtn-button' style={{display: 'none'}}/>
            <div className = "warehouse-mapp"></div>
            <div className = "info-message" id = "info-message"></div>
            
        </div>
    );
}

export default Settings;