import { useState } from "react";
import { useNavigate } from "react-router-dom";
import { useUser, UserType } from "@/contexts/UserContext";
import { Input } from "@/components/ui/input";
import { Button } from "@/components/ui/button";
import { ArrowLeft, User, Car } from "lucide-react";

const Register = () => {
  const [name, setName] = useState("");
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [userType, setUserType] = useState<UserType>("passenger");
  const { login } = useUser();
  const navigate = useNavigate();

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    if (!name || !email) return;
    login({ name, email, type: userType });
    navigate("/dashboard");
  };

  return (
    <div className="fixed inset-0 bg-card flex flex-col justify-center items-center z-50 p-6">
      <div className="w-full max-w-md">
        <button
          onClick={() => navigate("/")}
          className="mb-6 text-muted-foreground hover:text-foreground cursor-pointer flex items-center gap-2"
        >
          <ArrowLeft size={18} /> Voltar
        </button>

        <h2 className="text-3xl font-bold mb-6">Crie sua conta</h2>

        <form onSubmit={handleSubmit} className="space-y-4">
          <Input
            placeholder="Nome completo"
            value={name}
            onChange={(e) => setName(e.target.value)}
            required
            className="bg-secondary border-none focus-visible:ring-primary"
          />
          <Input
            type="email"
            placeholder="Email"
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

          <div className="flex gap-4">
            {(["passenger", "driver"] as const).map((type) => (
              <label key={type} className="flex-1 cursor-pointer group">
                <input
                  type="radio"
                  name="user-type"
                  value={type}
                  checked={userType === type}
                  onChange={() => setUserType(type)}
                  className="peer hidden"
                />
                <div className="p-4 border-2 border-border rounded-xl peer-checked:border-primary peer-checked:bg-secondary text-center transition group-hover:bg-secondary">
                  {type === "passenger" ? (
                    <User className="mx-auto mb-2" size={24} />
                  ) : (
                    <Car className="mx-auto mb-2" size={24} />
                  )}
                  <div className="font-semibold">
                    {type === "passenger" ? "Passageiro" : "Motorista"}
                  </div>
                </div>
              </label>
            ))}
          </div>

          <Button type="submit" className="w-full py-3 text-base font-bold rounded-lg">
            Cadastrar
          </Button>
        </form>
      </div>
    </div>
  );
};

export default Register;
