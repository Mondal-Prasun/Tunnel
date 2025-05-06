import { Link } from "react-router-dom";
import ContentImage from "../assets/demo.jpg";

function Contents({ item }) {
  return (
    <>
      <Link to={`/leech/${item.id}`} className="flex flex-col gap-2 items-center w-full h-full bg-white rounded shadow-md">
        <img src={ContentImage} alt={item.title} className="w-full rounded-t"/>
        <div className="w-full px-4 pb-4">
          <h2 className="font-bold">{item.title}</h2>
          <p>{item.description}</p>
        </div>
      </Link>
    </>
  );
}

export default Contents;
