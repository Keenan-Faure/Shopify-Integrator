import '../CSS/background.css';

/* Must start with a Caps letter */
function Background(props)
{
    return (
        <>
            <div className = "background" style = {{backgroundImage: `url(${props.image})`}}>
        
            </div>
            
        </>
    );
}

Background.defaultProps = 
{
    image: "#ccc"
};
export default Background;