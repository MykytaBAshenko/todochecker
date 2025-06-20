import { useForm } from "react-hook-form";
import { useNavigate } from "react-router-dom";
import { useAppDispatch, useAppSelector } from "../contexts/hooks";
import { registerUser } from "../contexts/actions/authActions";
import { useToast } from "../hooks/useToast";

export type RegisterFormData = {
  nickname: string;
  email: string;
  password: string;
  confirmPassword: string;
  avatar: string;
};

const Register = () => {
  const navigate = useNavigate();
  // We'll keep useAppContext for backward compatibility
  const dispatch = useAppDispatch();
    const { triggerToast } = useToast();
  
  // Get loading state from Redux
  const { isLoading, error } = useAppSelector(state => state.auth);

  const {
    register,
    watch,
    handleSubmit,
    formState: { errors },
  } = useForm<RegisterFormData>();

  const onSubmit = handleSubmit(async (data) => {
    // Add random avatar
    const formData = { 
      ...data,
      avatar: Array.from({ length: 8 }, () => Math.floor(Math.random() * 10)).join('')
    };
    
    const resultAction = await dispatch(registerUser(formData));
    
    if (registerUser.fulfilled.match(resultAction)) {
      // For backward compatibility
      triggerToast(({ message: "Registration Success!", type: "SUCCESS" }));
      navigate("/dashboard");
    }
  });

  return (
    <form className="flex flex-col gap-5" onSubmit={onSubmit}>
      <h2 className="text-3xl font-bold">Create an Account</h2>
        <label className="text-gray-700 text-sm font-bold flex-1">
          Nickname
          <input
            className="border rounded w-full py-1 px-2 font-normal"
            {...register("nickname", { required: "This field is required" })}
          ></input>
          {errors.nickname && (
            <span className="text-red-500">{errors.nickname.message}</span>
          )}
        </label>
      <label className="text-gray-700 text-sm font-bold flex-1">
        Email
        <input
          type="email"
          className="border rounded w-full py-1 px-2 font-normal"
          {...register("email", { required: "This field is required" })}
        ></input>
        {errors.email && (
          <span className="text-red-500">{errors.email.message}</span>
        )}
      </label>
      <label className="text-gray-700 text-sm font-bold flex-1">
        Password
        <input
          type="password"
          className="border rounded w-full py-1 px-2 font-normal"
          {...register("password", {
            required: "This field is required",
            minLength: {
              value: 6,
              message: "Password must be at least 6 characters",
            },
          })}
        ></input>
        {errors.password && (
          <span className="text-red-500">{errors.password.message}</span>
        )}
      </label>
      <label className="text-gray-700 text-sm font-bold flex-1">
        Confirm Password
        <input
          type="password"
          className="border rounded w-full py-1 px-2 font-normal"
          {...register("confirmPassword", {
            validate: (val) => {
              if (!val) {
                return "This field is required";
              } else if (watch("password") !== val) {
                return "Your passwords do no match";
              }
            },
          })}
        ></input>
        {errors.confirmPassword && (
          <span className="text-red-500">{errors.confirmPassword.message}</span>
        )}
      </label>
      <span>
        <button
          type="submit"
          className="bg-blue-600 text-white p-2 font-bold hover:bg-blue-500 text-xl"
          disabled={isLoading}
        >
          {isLoading ? "Creating Account..." : "Create Account"}
        </button>
      </span>
      {error && <p className="text-red-500">{error}</p>}
    </form>
  );
};

export default Register;