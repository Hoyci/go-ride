import { useState } from "react";
import { useNavigate } from "react-router-dom";
import { useUser } from "@/contexts/UserContext";
import { Input } from "@/components/ui/input";
import { Button } from "@/components/ui/button";
import { useMutation } from "@tanstack/react-query";
import { apiRequest } from "@/lib/api";
import { toast} from "sonner"

const Login = () => {
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const { login } = useUser();
  const navigate = useNavigate();

const loginMutation = useMutation({
  mutationFn: (credentials: any) => apiRequest("/login", "POST", credentials),
  onSuccess: (data) => {
    localStorage.setItem("access_token", data.access_token);
    localStorage.setItem("refresh_token", data.refresh_token);
    
    login({
      name: data.name,
      email: data.email,
      type: data.type.toLowerCase() as any,
    });
    navigate("/dashboard");
  },
  onError: (error: Error) => {
    toast.error("Credenciais inválidas ou erro de conexão.");
  }
});

const handleSubmit = (e: React.FormEvent) => {
  e.preventDefault();
  loginMutation.mutate({ email, password });
};

  return (
    <div className="fixed inset-0 bg-primary flex flex-col justify-center items-center z-50 p-6">
      <h1 className="text-4xl font-bold mb-8 text-primary-foreground">Uber Clone</h1>

      <div className="w-full max-w-md bg-card text-card-foreground p-6 rounded-2xl shadow-xl">
        <h2 className="text-2xl font-bold mb-6">Fazer Login</h2>

        <form onSubmit={handleSubmit} className="space-y-4">
          <Input
            type="email"
            placeholder="Email ou número de celular"
            value={email}
            onChange={(e) => setEmail(e.target.value)}
            required
            className="bg-secondary border-none focus-visible:ring-primary"
          />
          <Input
            type="password"
            placeholder="Senha"
            value={password}
            onChange={(e) => setPassword(e.target.value)}
            required
            className="bg-secondary border-none focus-visible:ring-primary"
          />
          <Button type="submit" className="w-full py-3 text-base font-bold rounded-lg">
            Entrar
          </Button>
        </form>

        <div className="mt-4 text-center text-sm text-muted-foreground">
          Não tem uma conta?{" "}
          <button
            onClick={() => navigate("/register")}
            className="text-foreground font-semibold underline cursor-pointer hover:text-muted-foreground"
          >
            Cadastre-se
          </button>
        </div>
      </div>
    </div>
  );
};

export default Login;
