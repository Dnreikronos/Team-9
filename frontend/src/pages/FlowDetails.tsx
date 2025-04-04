import { Link, useNavigate, useParams } from 'react-router';
import plusIcon from '../assets/plus.svg';
import TrainingDayContainer from '../components/TrainingDayContainer';
import { useEffect, useState } from 'react';
import getTrainingDayByFlowId from '../services/trainingDay/getByFlowId';
import Cookies from 'js-cookie';
import { Flow, TrainingDay } from '../types';
import getFlowById from '../services/flows/getById';
import { Trash } from 'phosphor-react';
import deleteFlow from '../services/flows/delete';

const FlowDetailsPage = () => {
    const [flow, setFlow] = useState<Flow>();
    const [daysOfTraining, setDaysOfTraining] = useState<TrainingDay[]>([]);

    const { id } = useParams();

    if (!id) return;

    const token = Cookies.get('auth_token');

    if (!token) throw new Error('JWT token invalid');

    const navigate = useNavigate();

    const handleClick = () => {
        navigate(`/add-new-day/${id}`);
    };

    const handleDeleteFlow = async () => {
        if (!flow?.id) return;

        const response = await deleteFlow(flow.id, token);

        if (response?.status !== 200) {
            throw new Error('Error to delete flow');
        }

        navigate('/');
    };

    useEffect(() => {
        const token = Cookies.get('auth_token');

        if (!token) throw new Error('JWT token invalid');

        const getFlow = async (flowId: string) => {
            const response = await getFlowById(flowId, token);

            if (response?.status === 200) {
                setFlow(response.data);
            }
        };

        const getDaysOfTraining = async (flowId: string) => {
            const response = await getTrainingDayByFlowId(flowId, token);

            if (response?.status === 200) {
                setDaysOfTraining(response.data);
            }
        };

        getDaysOfTraining(id);
        getFlow(id);
    }, []);

    return (
        <div className="flex w-full flex-col items-center p-7">
            <div className="flex flex-col items-center gap-4">
                <div className="mb-4 flex w-full flex-col gap-2 md:items-center">
                    <h1 className="text-2xl font-bold">{flow?.title}</h1>
                    <h3 className="text-xl text-gray-500">{flow?.level}</h3>
                </div>
                {daysOfTraining.length > 0 ? (
                    daysOfTraining.map((trainingDay) => (
                        <Link
                            key={trainingDay.id}
                            to={`/workout-info/${trainingDay.id}`}
                        >
                            <TrainingDayContainer
                                key={trainingDay.id}
                                trainingDay={trainingDay}
                            />
                        </Link>
                    ))
                ) : (
                    <span>No days of training added</span>
                )}
                <button
                    onClick={handleClick}
                    className="flex w-full cursor-pointer items-center justify-center gap-2.5 rounded-md bg-white p-3 text-xl text-black shadow-sm transition-colors duration-200 hover:bg-gradient-to-r hover:from-red-500 hover:to-purple-500 hover:text-white hover:transition hover:duration-500 focus:ring-2 focus:ring-gray-900 focus:ring-offset-2 focus:outline-none disabled:cursor-not-allowed disabled:opacity-50"
                >
                    <div className="flex items-center justify-center rounded-md bg-black p-3">
                        <img src={plusIcon} alt="plus" className="h-6 w-6" />
                    </div>
                    <span className="flex-grow text-lg font-medium">
                        Add training day
                    </span>
                </button>
                <button
                    onClick={handleDeleteFlow}
                    className="flex w-full cursor-pointer items-center justify-center gap-2.5 rounded-md bg-red-700 p-3 text-xl text-black shadow-sm transition-colors duration-200 hover:bg-red-800 hover:transition hover:duration-500 focus:ring-2 focus:ring-gray-900 focus:ring-offset-2 focus:outline-none disabled:cursor-not-allowed disabled:opacity-50"
                >
                    <div className="flex items-center justify-center rounded-md bg-black p-3">
                        <Trash color="white" className="h-6 w-6" />
                    </div>
                    <span className="flex-grow text-lg font-medium text-white">
                        Delete flow
                    </span>
                </button>
                <button
                    onClick={() => navigate('/')}
                    className="mt-4 w-full cursor-pointer rounded-md bg-gray-200 p-3 text-lg font-medium text-black transition-colors duration-200 hover:bg-gray-300 focus:ring-2 focus:ring-gray-900 focus:ring-offset-2 focus:outline-none"
                >
                    Back to Home
                </button>
            </div>
        </div>
    );
};

export default FlowDetailsPage;
