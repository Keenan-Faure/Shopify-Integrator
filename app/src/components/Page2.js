import {useEffect, useState} from 'react';
import { flushSync } from 'react-dom';
import { createRoot } from 'react-dom/client';
import Chart from "chart.js/auto";
import Queue_details from './semi-components/queue-details';
import Detailed_queue from './semi-components/Queue/detailed_queue';
import Background from './Background';
import $ from 'jquery';

import '../CSS/page1.css';


function Page2(props)
{
    useEffect(()=> 
    {
        /* Ensure the model is shown */
        let navbar = document.getElementById("navbar");
        navbar.style.display = "block";
        
        /* animation for the search bar */
        let search = document.querySelector(".search-area");
        setTimeout(() =>
        { search.style.opacity = "1"; search.style.animation = "appear 0.8s ease-in"; }, 1000);

        /* animation for the pan elements */
        let pan = document.querySelectorAll(".pan");
        let pag = document.getElementById("pag");
        setTimeout(() =>
        {
            for(let i = 0; i < pan.length; i ++)
            {
                pan[i].style.display = "block";
                pan[i].style.animation = "appear 0.8s ease-in";
            }
            pag.style.display = "block";
            pag.style.animation = "appear 1s ease-in";
        }, 1000);


        /* filter element animation */
        let filters = document.querySelector(".filter").children;
        setTimeout(() =>
        {
            for(let i = 0; i < filters.length; i ++)
            {
                filters[i].style.display = "block";
                filters[i].style.animation = "appear 0.8s ease-in";
            }
        }, 1000);

        /* filter image script to show when clicked on */
        let filter_button = document.getElementById("_filter");
        let filter = document.querySelectorAll(".filter-elements");
        let filter_img = document.querySelectorAll(".filter-img");
        let C_filter = document.getElementById("clear_filter");

        filter_button.disabled = true;
        C_filter.disabled = true;

        for(let i = 0; i < filter.length; i++)
        {
            /* Filter Onclick */
            filter[i].addEventListener("click", () =>
            {
                filter_img[i].style.display = "block";
                filter[i].style.backgroundColor = "rgba(64, 165, 24, 0.7)";

                filter_button.disabled = false;
                C_filter.disabled = false;
                filter_button.style.cursor = "pointer";
                C_filter.style.cursor = "pointer";

                for(let i = 0; i < filter.length; i++)
                {
                    filter[i].style.pointerEvents = "none";
                    filter[i].style.cursor = "not-allowed";
                } 
            });

            /* Clear Filter */
            C_filter.addEventListener("click", () =>
            {
                filter_img[i].style.display = "none";
                filter[i].style.backgroundColor = "rgba(61, 61, 61, 0.7)";

                filter_button.disabled = true;
                C_filter.disabled = true;
                filter_button.style.cursor = "not-allowed";
                C_filter.style.cursor = "not-allowed";

                for(let i = 0; i < filter.length; i++)
                {
                    filter[i].style.pointerEvents = "";
                    filter[i].style.cursor = "pointer";
                }
            });
        }

        filter_button.addEventListener("click", () =>
        {
            let next = document.getElementById("next");
            for(let i = 0; i < filter_img.length; i++)
            {
                if(filter_img[i].style.display == "block")
                {
                    let type = filter_img[i].nextSibling.className;
                    const api_key = localStorage.getItem('api_key');
                    $.ajaxSetup({ headers: { 'Authorization': 'ApiKey ' + api_key} });
                    $.get("http://localhost:8080/api/queue/filter?type=" + type, [], [], 'json')
                    .done(function( _data) 
                    {
                        console.log(_data);
                        if(_data.length < 10)
                        {
                            next.disabled = true;
                            next.style.cursor = "not-allowed";
                        }
                        if(_data == "")
                        {
                            document.querySelector(".pan-main").remove();
                            document.querySelector(".empty-message").style.display = "block";
                        }
                        else 
                        {
                            if(document.querySelector(".pan-main") != null)
                            {
                                document.querySelector(".pan-main").remove();
                                let pan_main = document.createElement('div');
                                let main_elements = document.querySelector(".main-elements");
                                pan_main.className = "pan-main";
                                main_elements.appendChild(pan_main);

                                let root = createRoot(pan_main);
                                flushSync(() => 
                                {
                                    root.render(_data.map((el, i) => <Queue_details key={`${el.title}_${i}`}
                                    Queue_Type={el.queue_type} Queue_Instruction={el.instruction} Queue_Status={el.status} Queue_ID={el.id}
                                    />))
                                });
                                setTimeout(() => { DetailedView();}, 200);
                            }
                            else 
                            {
                                let pan_main = document.createElement('div');
                                let main_elements = document.querySelector(".main-elements");
                                pan_main.className = "pan-main";
                                main_elements.appendChild(pan_main);

                                let root = createRoot(pan_main);
                                flushSync(() => 
                                {
                                    root.render(_data.map((el, i) => <Queue_details key={`${el.title}_${i}`}
                                    Queue_Type={el.queue_type} Queue_Instruction={el.instruction} Queue_Status={el.status} Queue_ID={el.id}
                                    />))
                                });
                                setTimeout(() => { DetailedView();}, 200);
                            }
                        }
                    })
                    .fail( function(xhr) { alert(xhr.responseText); });
                }
            }
            Filter_Pagintation(1);
        });

        function DetailedView()
        {
            let queue = document.querySelector(".queue");
            let pan = document.querySelectorAll(".pan");
            for(let i = 0; i < pan.length; i++)
            {
                pan[i].addEventListener("click", () =>
                {
                    let id = pan[i].querySelector(".p-d-id").innerHTML;
                    /*  API  */
                    const api_key = localStorage.getItem('api_key');
                    $.ajaxSetup({ headers: { 'Authorization': 'ApiKey ' + api_key} });
                    $.get("http://localhost:8080/api/queue/" + id, [], [], 'json')
                    .done(function(_data) 
                    {   
                        console.log(_data);
                        if(document.querySelector(".details") != null)
                        //div already exists, remove it, and create another
                        {

                            document.querySelector(".details").remove();
                            let details = document.createElement('div');
                            details.className = "details";
                            queue.appendChild(details);

                            let rot = createRoot(details);
                            rot.render( <Detailed_queue key={`${_data.title}_${i}`} Queue_Status={_data.status} Queue_Description={_data.description}
                            Queue_Type={_data.queue_type} Queue_Instruction={_data.instruction} Queue_ID={_data.id} Created_At={_data.created_at} 
                            Updated_At={_data.updated_at}
                            />)
                            
                        }
                        else 
                        //create new div
                        {
                            let details = document.createElement('details');
                            queue.appendChild(details);
                            let rot = createRoot(details);
                            rot.render( <Detailed_queue key={`${_data.title}_${i}`} Queue_Status={_data.status} Queue_Description={_data.description}
                            Queue_Type={_data.queue_type} Queue_Instruction={_data.instruction} Queue_ID={_data.id} Created_At={_data.created_at} 
                            Updated_At={_data.updated_at}
                            />)
                        }
                    })
                    .fail( function(xhr) 
                    {
                        alert(xhr.responseText);
                    });
                    setTimeout(() =>
                    {
                        let filter = document.querySelector(".filter");
                        let main = document.querySelector(".main");
                        let navbar = document.getElementById("navbar");
                        let details = document.querySelector(".details");

                        filter.style.animation = "Fadeout 0.5s ease-out";
                        main.style.animation = "Fadeout 0.5s ease-out";
                        navbar.style.animation = "Fadeout 0.5s ease-out";
                        filter.style.display = "none";
                        main.style.display = "none";
                        navbar.style.display = "none";
                        details.style.display = "block";
                    }, 50);
                });
            } 
        }

        /* Script to automatically format the number of elements on each page */
        const content = document.querySelector('.center'); 
        const paginationContainer = document.createElement('div');
        const paginationDiv = document.body.appendChild(paginationContainer);
        paginationContainer.classList.add('pagination');
        content.appendChild(paginationContainer);

        function Filter_Pagintation(index)
        {
            if(index == 1)
            {

                let next = document.getElementById("next");
                let ahead = index + 1;
                for(let i = 0; i < filter_img.length; i++)
                {
                    if(filter_img[i].style.display == "block")
                    {
                        let type = filter_img[i].nextSibling.className;
                        const api_key = localStorage.getItem('api_key');
                        $.ajaxSetup({ headers: { 'Authorization': 'ApiKey ' + api_key} });
                        $.get("http://localhost:8080/api/queue/filter?type=" + type + "?page="+ ahead, [], [], 'json')
                        .done(function( _data) 
                        {
                            console.log(_data);
                            if(_data.length < 10)
                            {
                                next.disabled = true;
                                next.style.cursor = "not-allowed";
                            }
                            if(_data == "") { let next = document.getElementById("next"); next.style.cursor = "not-allowed"; next.disabled = true; } 
                            else 
                            {
                                if(document.querySelector(".pan-main") != null)
                                {
                                    document.querySelector(".pan-main").remove();
                                    let pan_main = document.createElement('div');
                                    let main_elements = document.querySelector(".main-elements");
                                    pan_main.className = "pan-main";
                                    main_elements.appendChild(pan_main);

                                    let root = createRoot(pan_main);
                                    flushSync(() => 
                                    {
                                        root.render(_data.map((el, i) => <Queue_details key={`${el.title}_${i}`} 
                                        />))
                                    });
                                    setTimeout(() => { DetailedView();}, 200);
                                }
                                else 
                                {
                                    let pan_main = document.createElement('div');
                                    let main_elements = document.querySelector(".main-elements");
                                    pan_main.className = "pan-main";
                                    main_elements.appendChild(pan_main);

                                    let root = createRoot(pan_main);
                                    flushSync(() => 
                                    {
                                        root.render(_data.map((el, i) => <Queue_details key={`${el.title}_${i}`} 
                                        Queue_Type={el.queue_type} Queue_Instruction={el.instruction} Queue_Status={el.status} Queue_ID={el.id}
                                        />))
                                    });
                                    setTimeout(() => { DetailedView();}, 200);
                                }
                            }
                        })
                        .fail( function(xhr) { alert(xhr.responseText); });
                    }
                }
                
            }

            /* Check done to remove old elements if they exist */
            if(document.getElementById("next") != null && document.getElementById("prev") != null && document.getElementById("hod") != null)
            //If they exist remove them, and create new based on the new index value
            {
                document.getElementById("next").remove();
                document.getElementById("prev").remove();
                document.getElementById("hod").remove();

                const pageButton = document.createElement('button');
                pageButton.id = "hod";
                pageButton.className = "active";
                pageButton.innerHTML = index;
                paginationDiv.appendChild(pageButton);

                const nextPage = document.createElement('button');
                nextPage.id = "next";
                nextPage.innerHTML = "→";
                paginationDiv.appendChild(nextPage);

                const prevPage = document.createElement('button');
                prevPage.id = "prev";
                prevPage.innerHTML = "←";
                paginationDiv.appendChild(prevPage);
                if(index == 1) { prevPage.disabled = true; prevPage.style.cursor = "not-allowed"; }
                else if(index > 1) { prevPage.style.cursor = "pointer"; prevPage.disabled = false; nextPage.disabled = false; }
                else if(index <= 1) {prevPage.disabled = true; prevPage.style.cursor = "not-allowed"; }

                nextPage.addEventListener("click", () =>
                {
                    let category = document.querySelector(".category").innerHTML;
                    let type = document.querySelector(".type").innerHTML;
                    let vendor = document.querySelector(".vendor").innerHTML;

                    index = index + 1;
                    /*  API  */
                    const api_key = localStorage.getItem('api_key');
                    $.ajaxSetup({ headers: { 'Authorization': 'ApiKey ' + api_key} });
                    $.get("http://localhost:8080/api/products/filter?type=" + type + "&vendor=" + vendor + "&category=" + category + "&page=" + index, [], [])
                    .done(function( _data) 
                    {
                        console.log(_data);

                        document.querySelector(".pan-main").remove();
                        let div = document.createElement("div");
                        div.className = "pan-main";
                        div.id = "pan-main";
                        let main = document.querySelector(".main-elements");
                        main.appendChild(div);
                        let root = createRoot(div);
                        flushSync(() => 
                        { 
                            root.render(_data.map((el, i) => <Queue_details key={`${el.title}_${i}`}
                            Queue_Type={el.queue_type} Queue_Instruction={el.instruction} Queue_Status={el.status} Queue_ID={el.id}
                            />))
                        });
                        
                    })
                    .fail( function(xhr) { alert(xhr.responseText); });

                    let ahead = index + 1;
                    /*  API  */
                    $.get("http://localhost:8080/api/products/filter?type=" + type + "&vendor=" + vendor + "&category=" + category + "&page=" + ahead, [], [])
                    .done(function( _data) 
                    {
                        console.log(_data);
                        if(_data == "") { let next = document.getElementById("next"); next.style.cursor = "not-allowed"; next.disabled = true; }  
                    })
                    .fail( function(xhr) { alert(xhr.responseText); });

                    Filter_Pagintation(index);
                    setTimeout(() => { DetailedView();}, 200);
                });

                prevPage.addEventListener("click", () =>
                {
                    let category = document.querySelector(".category").innerHTML;
                    let type = document.querySelector(".type").innerHTML;
                    let vendor = document.querySelector(".vendor").innerHTML;
                    index = index - 1;
                    const api_key = localStorage.getItem('api_key');
                    $.ajaxSetup({ headers: { 'Authorization': 'ApiKey ' + api_key} });
                    $.get("http://localhost:8080/api/products/filter?type=" + type + "&vendor=" + vendor + "&category=" + category + "&page=" + index, [], [])
                    .done(function( _data) 
                    {
                        console.log(_data);
                        document.querySelector(".pan-main").remove();
                        let div = document.createElement("div");
                        div.className = "pan-main";
                        div.id = "pan-main";
                        let main = document.querySelector(".main-elements");
                        main.appendChild(div);
                        let root = createRoot(div);
                        flushSync(() => 
                        { 
                            root.render(_data.map((el, i) => <Queue_details key={`${el.title}_${i}`}
                            Queue_Type={el.queue_type} Queue_Instruction={el.instruction} Queue_Status={el.status} Queue_ID={el.id}
                            />))
                        });

                    
                    })
                    .fail( function(xhr) { alert(xhr.responseText); });
                    Filter_Pagintation(index);
                    setTimeout(() => { DetailedView();}, 200);
                });
            }
            else 
            //If they dont exist create new ones 
            {
                const pageButton = document.createElement('button');
                pageButton.id = "hod";
                pageButton.className = "active";
                pageButton.innerHTML = index;
                paginationDiv.appendChild(pageButton);

                const nextPage = document.createElement('button');
                nextPage.id = "next";
                nextPage.innerHTML = "→";
                paginationDiv.appendChild(nextPage);

                const prevPage = document.createElement('button');
                prevPage.id = "prev";
                prevPage.innerHTML = "←";
                paginationDiv.appendChild(prevPage);

                if(index == 1) { prevPage.disabled = true; prevPage.style.cursor = "not-allowed"; }
                else if(index > 1) { prevPage.style.cursor = "pointer"; prevPage.disabled = false; nextPage.disabled = false; }
                else if(index <= 1) {prevPage.disabled = true; prevPage.style.cursor = "not-allowed"; }

                nextPage.addEventListener("click", () =>
                {
                    let category = document.querySelector(".category").innerHTML;
                    let type = document.querySelector(".type").innerHTML;
                    let vendor = document.querySelector(".vendor").innerHTML;

                    index = index + 1;

                    const api_key = localStorage.getItem('api_key');
                    $.ajaxSetup({ headers: { 'Authorization': 'ApiKey ' + api_key} });
                    $.get("http://localhost:8080/api/products/filter?type=" + type + "&vendor=" + vendor + "&category=" + category + "&page=" + index, [], [])
                    .done(function( _data) 
                    {
                        console.log(_data);

                        document.querySelector(".pan-main").remove();
                        let div = document.createElement("div");
                        div.className = "pan-main";
                        div.id = "pan-main";
                        let main = document.querySelector(".main-elements");
                        main.appendChild(div);
                        let root = createRoot(div);
                        flushSync(() => 
                        { 
                            root.render(_data.map((el, i) => <Queue_details key={`${el.title}_${i}`} Queue_Updated_At={el.updated_at} Queue_Creation_Date={el.created_at} 
                            Queue_Type={el.queue_type} Queue_Instruction={el.instruction} Queue_Status={el.status} Queue_ID={el.id}
                            />)) 
                        });

                    })
                    .fail( function(xhr) { alert(xhr.responseText); });
                    Filter_Pagintation(index);
                    setTimeout(() => { DetailedView();}, 200);
                });

                prevPage.addEventListener("click", () =>
                {
                    let category = document.querySelector(".category").innerHTML;
                    let type = document.querySelector(".type").innerHTML;
                    let vendor = document.querySelector(".vendor").innerHTML;
                    index = index - 1;

                    const api_key = localStorage.getItem('api_key');
                    $.ajaxSetup({ headers: { 'Authorization': 'ApiKey ' + api_key} });
                    $.get("http://localhost:8080/api/products/filter?type=" + type + "&vendor=" + vendor + "&category=" + category + "&page=" + index, [], [])
                    .done(function( _data) 
                    {
                        console.log(_data);

                        document.querySelector(".pan-main").remove();
                        let div = document.createElement("div");
                        div.className = "pan-main";
                        div.id = "pan-main";
                        let main = document.querySelector(".main-elements");
                        main.appendChild(div);
                        let root = createRoot(div);
                        flushSync(() => 
                        { 
                            root.render(_data.map((el, i) => <Queue_details key={`${el.title}_${i}`}
                            Queue_Type={el.queue_type} Queue_Instruction={el.instruction} Queue_Status={el.status} Queue_ID={el.id}
                            />))
                        });
                        
                    })
                    .fail( function(xhr) { alert(xhr.responseText); });
                    Filter_Pagintation(index);
                    setTimeout(() => { DetailedView();}, 200);
                });
            } 
        }

        /* Queue Information Requests */
        let graph_data = {};
        const api_key = localStorage.getItem('api_key');
        $.ajaxSetup({ headers: { 'Authorization': 'ApiKey ' + api_key} });
        $.get("http://localhost:8080/api/queue/view", [], [], 'json')
        .done(function(_data) 
        {
            console.log(_data);
            graph_data = _data;
        })
        .fail( function(xhr) 
        {
            alert(xhr.responseText);
        });

        let ctx = document.getElementById("queue-graph");
        setTimeout(() =>
        {
            
            // FETCH GRAPH
            new Chart(ctx, {
                type: 'bar',
                data: {
                    labels: Object.keys(graph_data),
                    datasets: [{
                        label: 'Queue View',
                        data: Object.values(graph_data),
                        borderWidth: 1,
                        borderColor: "rgb(70 201 159)",
                        backgroundColor: "rgb(70 201 159)",
                        pointRadius: "5",
                    }]
                },
                options: {
                    scales: {
                        y: 
                        {
                            beginAtZero: true
                        }
                    }
                }
            });
        }, 200);

        let graph = document.querySelector(".graph");
        setTimeout(() =>
        {
            graph.style.animation = "SlideUp2 0.8s ease-in";
            graph.style.display = "block";
        }, 1200);
    
    }, []);

    return (
        <>
            <Background />
            <div className = "filter" style = {{display: props.filter_display}}>
                <div className = "filter-title"><b>Available Filters:</b></div>
                <div className = "filter-elements">
                    <div className = "filter-elements-text">Filter By Order</div>
                    <div className = "filter-img"/>
                    <div className = "order"></div>
                </div>
                <div className = "filter-elements" disabled>
                    <div className = "filter-elements-text">Filter By Product</div>
                    <div className = "filter-img"/>
                    <div className = "product"></div>
                </div>
                <div className = "filter-elements" disabled>
                    <div className = "filter-elements-text">Filter By Customer</div>
                    <div className = "filter-img"/>
                    <div className = "customer"></div>
                </div>
                <br />
                <div id = "button-hold">
                    <button id = "clear_filter"className = "filter-button">Clear Filter</button>
                    <button id = "_filter"className = "filter-button">Filter Results</button>
                </div>
                
                <br /><br /><br /><br /><br />
                <div className = "queue-view">
                    <div className="graph g1" style = {{width: '95%'}}>
                        <canvas id="queue-graph"></canvas>
                    </div>
                </div>
            </div>
        </>
    );
}
Page2.defaultProps = 
{
    filter_display: 'block', 
    main_display: 'block',
    main_bgc: '',
    main_top: '13%',
    main_left: '51%',
    transform: 'translate(-30%, -6%)',
    width: '70%',
    height: '96%', 
    animation: 'SlideUp2 0.8s ease-in'
};

export default Page2;