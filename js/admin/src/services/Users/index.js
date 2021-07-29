import { format as formatUrl } from 'url';

const createUserService = ({ httpClient }) => {
    const fetchUsers = async (filters = {}) => {
        const { search } = filters;
        const url = formatUrl({
            pathname: `${httpClient.getApiUrl()}/users`,
            query: {
                search,
            },
        });
        const { data } = await httpClient.get(url);

        return data;
    };

    const postUser = async (user) => {
        const url = `${httpClient.getApiUrl()}/users`;
        const { data } = await httpClient.post(url, user);

        return data;
    };

    const alterUser = async (id, changes) => {
        const url = `${httpClient.getApiUrl()}/users/${id}`;
        const { data } = await httpClient.put(url, changes);

        return data;
    };

    const fetchUserById = async (id) => {
        const url = `${httpClient.getApiUrl()}/users/${id}`;
        const { data } = await httpClient.get(url);

        return data;
    };

    const deleteUserById = async (id) => {
        const url = `${httpClient.getApiUrl()}/users/${id}`;
        const { data } = await httpClient.delete(url);

        return data;
    };

    return Object.freeze({
        fetchUsers,
        postUser,
        alterUser,
        fetchUserById,
        deleteUserById,
    });
};

export default createUserService;
